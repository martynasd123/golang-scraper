import * as React from "react"
import {useContext, useState} from "react"
import "./loginPage.less"
import {sendAuthRequest} from "../api/auth";
import {AuthContext} from "../contexts/AuthenticationContext";
import {useNavigate} from "react-router-dom";
import {AlertContext, AlertType} from "../contexts/AlertContext";
import {AxiosError, AxiosResponse} from "axios";

interface LoginFormState {
    username: string;
    password: string;
}

const LoginPage: React.FC = () => {
    const [credentials, setCredentials] = useState<LoginFormState>({
        username: "",
        password: ""
    });

    const { showAlert } = useContext(AlertContext)

    const {setUsername} = useContext(AuthContext)
    const navigate = useNavigate()

    const mapErrorMessage = (err: AxiosError) => {
        if(err.response) {
            if (err.response.status == 403) {
                return "Bad credentials"
            } else {
                return `Request failed: ${err.response.status}: ${err.response.data}`
            }
        } else {
            return `Unexpected error: ${err.message}.`
        }
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        sendAuthRequest(credentials)
            .then(response => {
                setUsername(credentials.username)
                navigate("/", {replace: true})
            }).catch(err => {
            showAlert({
                type: AlertType.WARNING,
                message: mapErrorMessage(err)
            })
        })
    };

    return <div className="login-page-root">
        <form onSubmit={handleSubmit} className="login-modal">
            <h2>Web scraping log-in</h2>
            <div className="login-form-group">
                <label htmlFor="username-input">Username</label>
                <input value={credentials.username}
                       type="text"
                       id="username-input"
                       onChange={e => setCredentials(prev => ({...prev, username: e.target.value}))}
                />
            </div>
            <div className="login-form-group">
                <label htmlFor="password-input">Password</label>
                <input value={credentials.password}
                       type="password"
                       id="password-input"
                       onChange={e => setCredentials(prev => ({...prev, password: e.target.value}))}
                />
            </div>
            <input type="submit" disabled={!credentials.username || !credentials.password} value="Log in"/>
        </form>
    </div>
}

export default LoginPage