import {AuthenticatedClient, Client} from "./client";
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

export enum TaskStatus {
    STATUS_PENDING = "PENDING",
    STATUS_INITIATING = "INITIATING",
    STATUS_TRYING_LINKS = "TRYING_LINKS",
    STATUS_FINISHED = "FINISHED",
    STATUS_ERROR = "ERROR",
    STATUS_INTERRUPTED = "INTERRUPTED",
    STATUS_INTERRUPTING = "INTERRUPTING"
}

export interface TaskStateUpdate {
    id?: number;
    status: TaskStatus;
    link: string;
    externalLinks?: number;
    internalLinks?: number;
    inaccessibleLinks?: number;
    htmlVersion?: string;
    pageTitle?: string;
    headingsByLevel?: [number, number, number, number, number, number];
    crawledLinks: number;
    error?: string;
}

export const sendAddTaskRequest = (request: AddTaskRequest) => AuthenticatedClient
    .post<AddTaskRequest, AxiosResponse<AddTaskResponse>>("/api/scrape/add-task", request)

export const sendInterruptTaskRequest = (id: string) => AuthenticatedClient
    .post<void, AxiosResponse<void>>(`/api/scrape/task/${id}/interrupt`)

export const sendGetTasksRequest = () => AuthenticatedClient
    .get<AddTaskRequest, AxiosResponse<TaskListItemResponse[]>>("/api/scrape/tasks")
