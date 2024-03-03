import { API } from "../backend/api";
import {
    useMutation,
    UseMutationOptions,
    useQuery,
} from "@tanstack/react-query";
import { Email } from "../backend/models";

export const useCurrentUserEmails = () => {
    const query = useQuery({
        queryKey: ["user_emails"],
        queryFn: API.getEmailsCurrentUser,
    });
    const {
        data: emails,
        isLoading: isLoadingEmails,
        error: errorEmails,
    } = query;
    return { emails, isLoadingEmails, errorEmails };
};

export const useCreateCurrentUserEmail = (
    options: UseMutationOptions<unknown, unknown, Partial<Email>>
) => {
    const {
        mutate: createEmail,
        isPending: isCreatingEmail,
        error: errorCreateEmail,
        reset: resetCreateEmail,
    } = useMutation({
        mutationKey: ["user_emails"],
        mutationFn: API.postEmailCurrentUser,
        ...options,
    });
    return {
        createEmail,
        isCreatingEmail,
        errorCreateEmail,
        resetCreateEmail,
    };
};

export const useDeleteCurrentUserEmail = (
    options: UseMutationOptions<unknown, unknown, Partial<Email>>
) => {
    const {
        mutate: deleteEmail,
        isPending: isDeletingEmail,
        error: errorDeleteEmail,
        reset: resetDeleteEmail,
    } = useMutation({
        mutationKey: ["user_emails"],
        mutationFn: API.deleteEmailCurrentUser,
        ...options,
    });
    return { deleteEmail, isDeletingEmail, errorDeleteEmail, resetDeleteEmail };
};
