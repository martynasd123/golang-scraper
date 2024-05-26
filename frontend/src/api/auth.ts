import {AuthenticatedClient, Client} from "./client";
import {AxiosResponse} from "axios";

interface AuthRequest {
    username: string;
    password: string;
}

interface LogOutRequest {
    username: string;
}

export const sendAuthRequest = (request: AuthRequest) => Client
    .post<AuthRequest, AxiosResponse<void>>("/api/auth/", request)

export const sendLogOutRequest = (request: LogOutRequest) => AuthenticatedClient
    .post<LogOutRequest, AxiosResponse<void>>("/api/auth/log-out", request)