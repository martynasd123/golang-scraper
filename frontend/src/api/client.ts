import axios from "axios";

export const Client = axios.create();

export const AuthenticatedClient = axios.create();

AuthenticatedClient.interceptors.response.use(null, async (err) => {
    if (err.response?.status === 403) {
        // Request was rejected - try to refresh token
        try {
            await axios({
                method: 'post',
                url: "/api/auth/refresh-token",
                data: {
                    username: localStorage.getItem("username")
                }
            });
            // Refresh token succeeded - retry original request
            axios(err.config);
        } catch {
            localStorage.removeItem("username")
            return await Promise.reject(err);
        }
    }
    return Promise.reject(err)
})
