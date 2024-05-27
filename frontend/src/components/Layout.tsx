import * as React from "react"
import {Outlet, useNavigate} from "react-router-dom";
import "./layout.less"
import {AxiosError} from "axios";
import {useContext} from "react";
import {AuthContext} from "../contexts/AuthenticationContext";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import {sendLogOutRequest} from "../api/auth";
import {IoMdArrowBack as BackArrow} from "react-icons/io";
import classNames from "../util/classNames";

function mapErrorMessage(err: AxiosError) {
    return `Unexpected error: ${err.message}.`
}

interface HeaderProps {
    backButtonLink?: string | null;
    progressBarProgress?: number | null;
    progressBarErroneous?: boolean;
}

function HeaderComponent({
                             progressBarProgress = null,
                             backButtonLink = null,
                             progressBarErroneous = null
                         }: HeaderProps) {
    const {username, setUsername} = useContext(AuthContext)
    const navigate = useNavigate()

    const {showAlert} = useContext(AlertContext)

    const handleLogOut = () => {
        sendLogOutRequest({username})
            .then(() => {
                setUsername(null)
                navigate("/login")
            }).catch(err => {
            showAlert({
                type: AlertType.WARNING,
                message: mapErrorMessage(err)
            })
        })
    };

    return <div className="header">
        <div className="header-content">
            <span>
                {backButtonLink != null && <div onClick={() => navigate(backButtonLink)} className="back-button">
                    <BackArrow/>
                </div>}
                <h2>
                Logged in as {username}
                </h2>
            </span>
            <input type="button" value="Log out" onClick={handleLogOut}/>
        </div>
        {progressBarProgress != null &&
            <div className="progress-indicator-container">
                <div className={classNames("progress-indicator", {"error": progressBarErroneous})}
                     style={{width: `${progressBarProgress}%`}}/>
            </div>}
    </div>
}

function PageContent({children}: React.PropsWithChildren) {
    return <div className="page-content-wrapper">
        <div className="page-content">
            {children}
        </div>
    </div>
}

function Layout() {
    return <div className="layout">
        <Outlet/>
    </div>
}

Layout.Content = PageContent
Layout.Header = HeaderComponent

export default Layout