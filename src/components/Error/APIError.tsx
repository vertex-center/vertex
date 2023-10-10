import ErrorBox from "./ErrorBox";
import { HTMLProps } from "react";
import { AxiosError } from "axios";

type APIErrorProps = HTMLProps<HTMLDivElement> & {
    error?: AxiosError;
};

export function APIError(props: Readonly<APIErrorProps>) {
    const { error, ...others } = props;

    if (!error) return null;

    if (error.response)
        return <ErrorBox error={error.response.data} {...others} />;

    if (error.request)
        return <ErrorBox error="No response from server." {...others} />;

    return <ErrorBox error={error?.message} {...others} />;
}
