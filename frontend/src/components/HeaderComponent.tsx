import * as React from "react"
import {useContext} from "react"
import "./header.less"
import {AuthContext} from "../contexts/AuthenticationContext";
import {sendLogOutRequest} from "../api/auth";
import {useNavigate} from "react-router-dom";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import {AxiosError} from "axios";

function mapErrorMessage(err: AxiosError) {
    return `Unexpected error: ${err.message}.`
}

const Layout: React.FC = () => {
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
        <h2>
            Logged in as {username}
        </h2>
        <input type="button" value="Log out" onClick={handleLogOut}/>
    </div>
}

export default Layout