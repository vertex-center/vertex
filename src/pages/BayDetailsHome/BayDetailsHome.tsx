import { Fragment, useEffect } from "react";
import { Title } from "../../components/Text/Text";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsHome.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import classNames from "classnames";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import useInstance from "../../hooks/useInstance";

export default function BayDetailsHome() {
    const { uuid } = useParams();

    const { instance, setInstance } = useInstance(uuid);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(`/instance/${uuid}/events`);

        const onStatusChange = (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [uuid]);

    return (
        <Fragment>
            <Title className={styles.title}>Home</Title>
            <nav className={styles.nav}>
                {instance?.urls &&
                    instance?.urls
                        .filter((u) => u.kind === "client")
                        .map((u) => {
                            const portEnvName = instance?.environment
                                ?.filter((e) => e.type === "port")
                                ?.find((e) => e.default === u.port).name;

                            const port = instance?.env[portEnvName] ?? u.port;
                            const disabled = instance.status !== "running";

                            // @ts-ignore
                            let url = new URL(window.apiURL);
                            url.port = port;

                            return (
                                <a
                                    href={url.href}
                                    target="_blank"
                                    rel="noreferrer"
                                    className={classNames({
                                        [styles.navItem]: true,
                                        [styles.navItemDisabled]: disabled,
                                    })}
                                >
                                    <Horizontal>
                                        <Symbol
                                            className={styles.navItemSymbol}
                                            name="public"
                                        />
                                        <Spacer />
                                        <Symbol name="open_in_new" />
                                    </Horizontal>
                                    <div>{url.href}</div>
                                </a>
                            );
                        })}
            </nav>
        </Fragment>
    );
}
