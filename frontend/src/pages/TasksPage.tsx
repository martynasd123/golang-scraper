import * as React from "react"
import {useContext, useEffect, useState} from "react"
import CardComponent from "../components/CardComponent";
import "./tasksPage.less"
import {sendAddTaskRequest, sendGetTasksRequest} from "../api/scrape";
import {useNavigate} from "react-router-dom";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import {AxiosError} from "axios";
import { FaAngleRight as RightArrow } from "react-icons/fa6";
import Layout from "../components/Layout";

interface TaskCard {
    id: number;
    link: string;
    inaccessibleLinks?: number;
    pageTitle?: string;
    crawledLinks?: number;
    error?: string;
    status: string;
}

const mapAddTaskErrorMessage = (err: AxiosError) => {
    if (err.response) {
        if (err.response.status == 400 && err.response.data === "Invalid URL") {
            return "Invalid link provided"
        } else {
            return `Request failed: ${err.response.status} - ${err.response.data}`
        }
    } else {
        return `Unexpected error: ${err.message}.`
    }
};


const mapGetTasksErrorMessage = (err: AxiosError) => {
    if (err.response) {
        return `Unexpected error: ${err.message}.`
    }
};

const TasksPage: React.FC = () => {
    const [tasks, setTasks] = useState<TaskCard[]>()
    const [newScrapeTaskUrl, setNewScrapeTaskUrl] = useState("")

    const navigate = useNavigate()

    const {showAlert} = useContext(AlertContext)

    const handleAddNewTask = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        sendAddTaskRequest({link: newScrapeTaskUrl})
            .then(response => {
                navigate(`/task/${response.data.id}`)
            }).catch(err => showAlert({
            type: AlertType.WARNING,
            message: mapAddTaskErrorMessage(err)
        }))
    };

    useEffect(() => {
        sendGetTasksRequest().then(response => {
            setTasks(response.data)
        }).catch(err => showAlert({
            type: AlertType.WARNING,
            message: mapGetTasksErrorMessage(err)
        }))
    }, [])

    return <>
        <Layout.Header />
        <Layout.Content>
            <form className="new-task-form" onSubmit={handleAddNewTask}>
                <h2>Add a new scraping task</h2>
                <div className="input-group">
                    <label htmlFor="new-scrape-task-url-input">Enter URL to scrape</label>
                    <input value={newScrapeTaskUrl}
                           type="url"
                           id="new-scrape-task-url-input"
                           placeholder="www.your-url.com"
                           onChange={e => setNewScrapeTaskUrl(e.target.value)}
                    />
                </div>
                <input type="submit" value="Start scraping"/>
            </form>
            <h2>History</h2>
            {tasks && tasks.map(task => <CardComponent className="task-card" key={task.id}>
                <div>
                    <h3>{task.link}</h3>
                    {task.pageTitle && <div>
                        {task.pageTitle}
                    </div>}
                    {<div>Task status: {task.status}</div>}
                    {task.error && <div>
                        Error encountered: {task.error}
                    </div>}
                </div>
                <div onClick={() => navigate(`/task/${task.id}`)}>
                    <RightArrow/>
                </div>
            </CardComponent>)}
        </Layout.Content>
    </>
}

export default TasksPage