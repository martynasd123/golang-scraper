import * as React from "react"
import {Outlet} from "react-router-dom";

const AuthenticatedBoundary: React.FC = () => {
    return <Outlet/>
}

export default AuthenticatedBoundary