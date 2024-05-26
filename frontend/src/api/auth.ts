import {Client} from "./client";
import {AxiosResponse} from "axios";

interface AuthRequest {
    username: string;
    password: string;
}

export const sendAuthRequest = (request: AuthRequest) => Client
    .post<AuthRequest, AxiosResponse<void>>("/api/auth/", request)