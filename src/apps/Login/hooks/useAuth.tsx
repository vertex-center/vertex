import { getAuthToken } from "../../../backend/api/backend";

export default function useAuth(uuid?: string) {
    const token = getAuthToken();
    return { token, isLoggedIn: !!token };
}
