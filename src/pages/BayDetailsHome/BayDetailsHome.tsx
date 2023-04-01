import { Fragment, useCallback, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./BayDetailsHome.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import { getInstance, Instance, route } from "../../backend/backend";
import classNames from "classnames";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";

export default function BayDetailsHome() {
    const { uuid } = useParams();

    const [instance, setInstance] = useState<Instance>();

    const fetchInstance = useCallback(() => {
        getInstance(uuid).then((instance: Instance) => {
            setInstance(instance);
        });
    }, [uuid]);

    useEffect(() => {
        fetchInstance();
    }, [fetchInstance, uuid]);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(route(`/instance/${uuid}/events`));

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
            <Title>Home</Title>
            <nav className={styles.nav}>
                {instance?.environment
                    .filter((e) => e.type === "port")
                    .map((e) => {
                        console.log(instance.env);
                        const port = instance.env[e.name] ?? e.default;
                        const disabled = instance.status !== "running";
                        const href = `http://localhost:${port}`;

                        return (
                            <a
                                href={href}
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
                                <div>http://localhost:{port}</div>
                            </a>
                        );
                    })}
            </nav>
        </Fragment>
    );
}
