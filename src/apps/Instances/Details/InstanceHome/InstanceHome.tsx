import { Fragment } from "react";
import { Title } from "../../../../components/Text/Text";
import Symbol from "../../../../components/Symbol/Symbol";

import styles from "./InstanceHome.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import classNames from "classnames";
import useInstance from "../../../../hooks/useInstance";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function InstanceHome() {
    const { uuid } = useParams();

    const { instance, loading } = useInstance(uuid);

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <Title className={styles.title}>URLs</Title>
            <nav className={styles.nav}>
                {instance?.service?.urls &&
                    instance?.service?.urls
                        .filter((u) => u.kind === "client")
                        .map((u) => {
                            const portEnvName = instance?.service?.environment
                                ?.filter((e) => e.type === "port")
                                ?.find((e) => e.default === u.port).name;

                            const port =
                                instance?.environment[portEnvName] ?? u.port;
                            const disabled = instance.status !== "running";

                            // @ts-ignore
                            let url = new URL(window.apiURL);
                            url.port = port;
                            url.pathname = u.home ?? "";

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