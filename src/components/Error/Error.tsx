import styles from "./Error.module.sass";
import Symbol from "../Symbol/Symbol";
import { AxiosError } from "axios";
import { HTMLProps } from "react";
import classNames from "classnames";

export function Errors(props: HTMLProps<HTMLDivElement>) {
    const { children, className } = props;

    if (!children) return null;

    return <div {...props} className={classNames(styles.errors, className)} />;
}

type Props = HTMLProps<HTMLDivElement> & {
    error?: any;
};

export default function Error(props: Props) {
    const { error, className, ...others } = props;

    let err = error?.message ?? error;

    return (
        <div className={classNames(styles.box, className)} {...others}>
            <div className={styles.error}>
                <Symbol className={styles.symbol} name="error" />
                <h1>Error</h1>
            </div>
            <div className={styles.content}>
                {err ?? "An unknown error has occurred."}
            </div>
        </div>
    );
}

type APIErrorProps = HTMLProps<HTMLDivElement> & {
    error?: AxiosError;
};

export function APIError(props: APIErrorProps) {
    const { error, ...others } = props;

    if (!error) return null;

    if (error.response)
        return <Error error={error.response.data} {...others} />;

    if (error.request)
        return <Error error="No response from server." {...others} />;

    return <Error error={error.message} {...others} />;
}
