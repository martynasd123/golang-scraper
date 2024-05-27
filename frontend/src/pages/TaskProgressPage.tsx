import * as React from "react"
import {useContext, useEffect, useState} from "react"
import {useParams} from "react-router-dom";
import {sendInterruptTaskRequest, TaskStateUpdate, TaskStatus} from "../api/scrape";
import "./taskProgressPage.less"
import LoaderComponent from "../components/LoaderComponent";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import Layout from "../components/Layout";
import {AxiosError} from "axios";
import {FaRegCircleStop as StopIcon} from "react-icons/fa6";
import classNames from "../util/classNames";

const TaskInfoField: React.FC<{ name: string, value: string | number }> = ({name, value}) => {
    return <div className="task-info-field">
        <div className="row">
            <span>
                {name}
            </span>
            <span>
                {value == null ? "-" : value}
            </span>
        </div>
    </div>
}

const TaskInfoFields: React.FC<{ task: TaskStateUpdate }> = ({task}) => {
    return <div className="task-info-fields-container">
        <TaskInfoField name="Link" value={task.link}/>
        <TaskInfoField name="Page title" value={task.pageTitle}/>
        <TaskInfoField name="Internal links" value={task.internalLinks}/>
        <TaskInfoField name="External links" value={task.externalLinks}/>
        <TaskInfoField name="Inaccessible links" value={task.inaccessibleLinks}/>
        <TaskInfoField name="Contains login form" value={task.loginFormPresent == null ? null : task.loginFormPresent.toString()}/>
        <TaskInfoField name="Links crawled in total" value={task.crawledLinks}/>
        <TaskInfoField name="Html version" value={task.htmlVersion}/>
        {(task.headingsByLevel || new Array(6).fill(null)).map((number, i) =>
            <TaskInfoField key={i} name={`<h${i + 1}> headings`} value={number}/>)}
    </div>
}

function mapErrorMessage(err: AxiosError) {
    if (err.response?.status == 400) {
        if (err.response.data == "task already in final state") {
            return "Task already finished"
        } else if (err.response.data == "interrupt already sent") {
            return "Interrupt signal already sent"
        }
    }
    return "Unexpected error"
}

const TaskProgressPage: React.FC = () => {
    const [taskState, setTaskState] = useState<TaskStateUpdate>()
    const {id} = useParams()

    const [err, setErr] = useState<string>(null)

    const {showAlert} = useContext(AlertContext)

    const getProgressPercentage: () => number = () => {
        if (taskState == null) {
            return null
        }
        if (taskState.status == TaskStatus.STATUS_INITIATING || taskState.status == TaskStatus.STATUS_PENDING) {
            return 0
        }
        const totalLinks = taskState.externalLinks + taskState.internalLinks
        if (totalLinks == 0) {
            return 100
        }
        return ((taskState.crawledLinks * 100) / totalLinks)
    };

    useEffect(() => {
        const eventSource = new EventSource(`/api/scrape/task/${id}/listen`);

        eventSource.onmessage = (update) => {
            setTaskState(JSON.parse(update.data))
        }

        eventSource.onerror = (error) => {
            if (eventSource.readyState !== EventSource.CLOSED) {
                return
            }
            console.error(error)
            const msg = "Could not retrieve task info"
            showAlert({
                type: AlertType.WARNING,
                message: msg
            })
            setErr(msg)
            eventSource.close();
        }

        return () => {
            eventSource.close();
        };
    }, [])

    const isInterruptibleState = [TaskStatus.STATUS_PENDING,
        TaskStatus.STATUS_INITIATING,
        TaskStatus.STATUS_TRYING_LINKS].includes(taskState?.status)

    let PageContent;

    const interruptTask = () => {
        if (!isInterruptibleState) {
            return
        }
        sendInterruptTaskRequest(id)
            .then(() => {
                showAlert({
                    type: AlertType.SUCCESS,
                    message: "Interrupt signal sent"
                })
            })
            .catch((err) => {
                showAlert({
                    type: AlertType.WARNING,
                    message: mapErrorMessage(err)
                })
            })
    };
    if (err) {
        PageContent = <div>Error occurred while retrieving this task</div>
    } else if (!taskState) {
        PageContent = <LoaderComponent className="task-page-loader-wrapper"/>
    } else {
        PageContent = <>
            <div className="task-page-title-wrapper">
                <h2>Task #{taskState.id}</h2>
                <div onClick={interruptTask}
                     className={classNames("interrupt-icon", {disabled: !isInterruptibleState})}>
                    <StopIcon size={24}/>
                </div>
            </div>
            <div>Status: {taskState.status}
            </div>
            {taskState.error ?
                <div className="err-container">Encountered error: {taskState.error}</div>
                : <TaskInfoFields task={taskState}/>}
        </>
    }

    const isErrorState = taskState?.status === TaskStatus.STATUS_ERROR
        || taskState?.status == TaskStatus.STATUS_INTERRUPTED

    return <>
        <Layout.Header backButtonLink={"/"} progressBarErroneous={isErrorState}
                       progressBarProgress={getProgressPercentage()}/>
        <Layout.Content>
            {PageContent}
        </Layout.Content>
    </>
}

export default TaskProgressPage