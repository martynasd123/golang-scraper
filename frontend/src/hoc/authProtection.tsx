import * as React from "react"
import {useContext} from "react";
import {AuthContext} from "../contexts/AuthenticationContext";
import {Navigate} from "react-router-dom";

const authBoundary = (component: React.FC, requiresAuth: boolean, redirectPath: string): React.ReactNode => {
    const WithAuth = () => {
        const {isAuthenticated} = useContext(AuthContext);

        const Component = component
        if (isAuthenticated == requiresAuth) {
            return <Component/>;
        } else {
            return <Navigate to={redirectPath}/>;
        }
    }
    return <WithAuth/>
};

export const unauthenticatedOnly = (component: React.FC): React.ReactNode => authBoundary(component, false, "/");
export const authenticatedOnly = (component: React.FC): React.ReactNode => authBoundary(component, true, "/login");