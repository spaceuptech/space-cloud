package istio

import (
	"fmt"

	b64 "encoding/base64"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
	"github.com/spaceuptech/space-cloud/runner/model"
)

// CreateSecret is used to upsert secret
func (i *Istio) CreateSecret(projectID string, secretObj *model.Secret) error {
	// check whether the secret type is correct!
	if secretObj.Type != model.FileType && secretObj.Type != model.EnvType && secretObj.Type != model.DockerType {
		return fmt.Errorf("invalid secret type (%s) provided", secretObj.Type)
	}

	_, err := i.kube.CoreV1().Secrets(projectID).Get(secretObj.Name, metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// Create a new Secret
		logrus.Debugf("Creating secret (%s)", secretObj.Name)
		newSecret := generateSecret(projectID, secretObj)
		_, err := i.kube.CoreV1().Secrets(projectID).Create(newSecret)
		if err != nil {
			return err
		}

	} else if err == nil {
		// secret already exists...update it!
		logrus.Debugf("Updating secret (%s)", secretObj.Name)
		newSecret := generateSecret(projectID, secretObj)
		_, err = i.kube.CoreV1().Secrets(projectID).Update(newSecret)
		if err != nil {
			return err
		}
		return nil
	}
	logrus.Errorf("Failed to create secret (%s) - %s", secretObj.Name, err)
	return err
}

// ListSecrets lists all the secrets in the provided name-space!
func (i *Istio) ListSecrets(projectID string) ([]*model.Secret, error) {

	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorf("Failed to fetch list of secrets - %s", err)
		return nil, err
	}
	listOfSecrets := make([]*model.Secret, len(kubeSecret.Items))
	// Modifying SecretValue with empty []byte
	for i, v := range kubeSecret.Items {
		s := &model.Secret{
			Name:     v.ObjectMeta.Name,
			Type:     v.ObjectMeta.Annotations["secretType"],
			RootPath: v.ObjectMeta.Annotations["rootPath"],
			Data:     make(map[string]string, len(v.Data)),
		}
		for k1 := range v.Data {
			s.Data[k1] = ""
		}
		listOfSecrets[i] = s
	}
	return listOfSecrets, nil
}

// DeleteSecret is used to delete secrets!
func (i *Istio) DeleteSecret(projectID string, secretName string) error {
	err := i.kube.CoreV1().Secrets(projectID).Delete(secretName, &metav1.DeleteOptions{})
	if kubeErrors.IsNotFound(err) || err == nil {
		return nil
	}
	logrus.Errorf("Failed to delete secret (%s) - %s", secretName, err)
	return err
}

// SetKey adds a new secret key-value pair
func (i *Istio) SetKey(projectID string, secretName string, secretKey string, secretValObj *model.SecretValue) error {
	if secretName == "" || secretValObj.Value == "" {
		logrus.Errorf("Empty key/value provided. Key not set")
		return fmt.Errorf("key/value not provided; got (%s,%s)", secretName, secretValObj.Value)
	}
	// encoding secret value to base64
	encSecret := b64.StdEncoding.EncodeToString([]byte(secretValObj.Value))
	//Get secret and then check type
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return err
	} else if err == nil {
		//Add secret key-value
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			return fmt.Errorf("setKey operation cannot be performed on secrets with type docker")
		case v1.SecretTypeOpaque:
			kubeSecret.Data[secretKey] = []byte(encSecret)
		default:
			//Throw error
			return fmt.Errorf("invalid secret type - %s", kubeSecret.Type)
		}

		//Update the secret
		_, err := i.kube.CoreV1().Secrets(projectID).Update(kubeSecret)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteKey is used to delete a key from the secret!
func (i *Istio) DeleteKey(projectID string, secretName string, secretKey string) error {
	// Get secret
	kubeSecret, err := i.kube.CoreV1().Secrets(projectID).Get(secretName, metav1.GetOptions{})

	if kubeErrors.IsNotFound(err) {
		return fmt.Errorf("secret with name (%s) does not exist- %s", secretName, err)
	} else if err == nil {
		// check the type of secret (docker/opaque)
		switch kubeSecret.Type {
		case v1.SecretTypeDockerConfigJson:
			return fmt.Errorf("deleteKey operation cannot be performed on secrets with type docker")
		case v1.SecretTypeOpaque:
			delete(kubeSecret.Data, secretKey)
		default:
			// Throw error
			return fmt.Errorf("invalid secret type - %s", kubeSecret.Type)
		}
		// Update the secret
		_, err := i.kube.CoreV1().Secrets(projectID).Update(kubeSecret)
		if err != nil {
			return err
		}
	}
	return nil
}

// helper function
func generateSecret(projectID string, secret *model.Secret) *v1.Secret {
	encodedData := map[string][]byte{}
	var typeOfSecret v1.SecretType

	// Check what type of secret is to be created: file/env/docker
	if secret.Type == model.FileType || secret.Type == model.EnvType {
		typeOfSecret = v1.SecretTypeOpaque
		// Base64 encoding!
		for k, v := range secret.Data {
			encValue := b64.StdEncoding.EncodeToString([]byte(v))
			encodedData[k] = []byte(encValue)
		}
	} else if secret.Type == model.DockerType {
		typeOfSecret = v1.SecretTypeDockerConfigJson
		authSecret := secret.Data["username"] + ":" + secret.Data["password"]
		encAuthSecret := b64.StdEncoding.EncodeToString([]byte(authSecret))
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials
		dockerJSON := map[string]interface{}{
			"auths": map[string]interface{}{
				secret.Data["url"]: map[string]string{
					"auth": encAuthSecret,
				},
			},
		}
		data, _ := json.Marshal(dockerJSON)
		// encode the entire json object in base64 encoding
		encValue := b64.StdEncoding.EncodeToString([]byte(data))
		encodedData[v1.DockerConfigJsonKey] = []byte(encValue)
	}
	return &v1.Secret{Type: typeOfSecret, ObjectMeta: metav1.ObjectMeta{Name: secret.Name, Namespace: projectID, Annotations: map[string]string{"rootPath": secret.RootPath, "secretType": secret.Type}}, Data: encodedData}
}
