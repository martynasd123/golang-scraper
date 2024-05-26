import {Client} from "./client";
import {AxiosResponse} from "axios";

interface AddTaskRequest {
    link: string
}

interface AddTaskResponse {
    id: number
}

interface TaskListItemResponse {
    id: number;
    link: string;
    inaccessibleLinks?: number;
    pageTitle?: string;
    crawledLinks?: number;
    error?: string;
    status: string;
}

export const sendAddTaskRequest = (request: AddTaskRequest) => Client
    .post<AddTaskRequest, AxiosResponse<AddTaskResponse>>("/api/scrape/add-task", request)

export const sendGetTasksRequest = () => Client
    .get<AddTaskRequest, AxiosResponse<TaskListItemResponse[]>>("/api/scrape/tasks")