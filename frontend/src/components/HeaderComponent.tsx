import * as React from "react"
import {useContext} from "react"
import "./header.less"
import {AuthContext} from "../contexts/AuthenticationContext";
import {sendLogOutRequest} from "../api/auth";
import {useNavigate} from "react-router-dom";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import {AxiosError} from "axios";
import {HeaderContext} from "../contexts/HeaderContext";
import { IoMdArrowBack as BackArrow } from "react-icons/io";

function mapErrorMessage(err: AxiosError) {
    return `Unexpected error: ${err.message}.`
}

const Layout: React.FC = () => {
    const {progressBarProgress, backButtonLink} = useContext(HeaderContext)

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
                <div className="progress-indicator" style={{width: `${progressBarProgress}%`}}/>
            </div>}
    </div>
}

export default Layout