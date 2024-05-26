import * as React from "react"
import {useContext, useEffect, useState} from "react"
import {useParams} from "react-router-dom";
import {TaskStateUpdate, TaskStatus} from "../api/scrape";
import "./taskProgressPage.less"
import {HeaderContext} from "../contexts/HeaderContext";
import LoaderComponent from "../components/LoaderComponent";
import {AlertContext, AlertType} from "../contexts/AlertContext";

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
        <TaskInfoField name="Internal links" value={task.internalLinks}/>
        <TaskInfoField name="External links" value={task.externalLinks}/>
        <TaskInfoField name="Inaccessible links" value={task.inaccessibleLinks}/>
        <TaskInfoField name="Links crawled in total" value={task.crawledLinks}/>
        <TaskInfoField name="Html version" value={task.htmlVersion}/>
        {(task.headingsByLevel || new Array(6).fill(null)).map((number, i) =>
            <TaskInfoField key={i} name={`<h${i+1}> headings`} value={number}/>)}
    </div>
}

const TaskProgressPage: React.FC = () => {
    const [taskState, setTaskState] = useState<TaskStateUpdate>()
    const {id} = useParams()

    const { setProgressBarProgress, setBackButtonLink } = useContext(HeaderContext)
    const [err, setErr] = useState<string>(null)

    const {showAlert} = useContext(AlertContext)

    const getProgressPercentage = () => {
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
        setProgressBarProgress(getProgressPercentage())
        setBackButtonLink("/")
    }, [taskState])

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

    if (err) {
        return <div>Error occurred while retrieving this task</div>
    }

    if (!taskState) {
        return <LoaderComponent className="task-page-loader-wrapper"/>
    }

    return <>
        <h2>Task #{taskState.id}</h2>
        <div>Status: {taskState.status}</div>
        {taskState.error ? <div className="err-container">Encountered error: {taskState.error}</div>
            : <TaskInfoFields task={taskState}/>}

        <div className="progress-indicator" style={{
            width: `${getProgressPercentage()}%`
        }}></div>

    </>
}

export default TaskProgressPage