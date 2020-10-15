package com.spaceuptech.dbevents.spacecloud

import akka.actor.typed.{ActorRef, ActorSystem, Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import com.spaceuptech.dbevents.Global

import scala.concurrent.ExecutionContextExecutor
import scala.concurrent.duration._
import scala.util.{Failure, Success}

object ProjectsSupervisor {
  def apply(): Behavior[Command] = Behaviors.withTimers(timers => Behaviors.setup(context => new ProjectsSupervisor(context, timers)))

  sealed trait Command

  case class FetchProjects() extends Command
  case class Stop() extends Command

  val fetchProjectsKey: String = "fetch-projects"
}

class ProjectsSupervisor(context: ActorContext[ProjectsSupervisor.Command], timers: TimerScheduler[ProjectsSupervisor.Command]) extends AbstractBehavior[ProjectsSupervisor.Command](context) {

  import ProjectsSupervisor._

  // Member variables
  private var projectIdToActor: Map[String, ActorRef[ProjectManager.Command]] = Map.empty

  // Start the timer to fetch projects
  timers.startTimerAtFixedRate(fetchProjectsKey, FetchProjects(), 1.minute)

  override def onMessage(msg: Command): Behavior[Command] = {
    msg match {
      case FetchProjects() =>
        fetchProjects()
        this

      case Stop() => Behaviors.stopped
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      // Stop the timers
      timers.cancelAll()

      // Stop all children
      for ((_, actor) <- projectIdToActor) {
        actor ! ProjectManager.Stop()
      }
      projectIdToActor = Map.empty
      this

  }

  private def fetchProjects(): Unit = {
    implicit val system: ActorSystem[Nothing] = context.system
    implicit val executionContext: ExecutionContextExecutor = system.executionContext

    // Make http request
    fetchSpaceCloudResource[Project](s"http://${Global.gatewayUrl}/v1/config/projects/*").onComplete {
      case Success(res) => processProjects(res.result)
      case Failure(ex) => context.log.error("Unable to fetch projects", ex)
    }
  }

  private def processProjects(projects: Array[Project]): Unit = {
    // Create an actor for new projects
    projects.foreach(project => {
      if (!this.projectIdToActor.contains(project.id)) {
        val actor = context.spawn(ProjectManager(project.id), project.id)
        actor ! ProjectManager.FetchEventingConfig()
        projectIdToActor += project.id -> actor
      }
    })

    // Close old project actors
    this.projectIdToActor = this.projectIdToActor.filter(elem => removeProjectIfInactive(projects, elem._1, elem._2))
  }

  private def removeProjectIfInactive(projects: Array[Project], projectId: String, actor: ActorRef[ProjectManager.Command]): Boolean = {
    if (!projects.exists(project => project.id == projectId)) {
      actor ! ProjectManager.Stop()
      return false
    }
    true
  }
}
