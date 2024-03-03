import { getAuthToken } from "../../../backend/server";

export default function useAuth(uuid?: string) {
    const token = getAuthToken();
    return { token, isLoggedIn: !!token };
}
