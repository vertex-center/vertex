import { MaterialIcon, Title } from "@vertex-center/components";
import styles from "./ContainerHome.module.sass";
import { useParams } from "react-router-dom";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import classNames from "classnames";
import useContainer from "../../hooks/useContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { APIError } from "../../../../components/Error/APIError";
import Content from "../../../../components/Content/Content";

export default function ContainerHome() {
    const { uuid } = useParams();

    const { container, isLoading, error } = useContainer(uuid);

    return (
        <Content>
            <Title variant="h2">URLs</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <nav className={styles.nav}>
                {container?.service?.urls &&
                    container?.service?.urls
                        .filter((u) => u.kind === "client")
                        .map((u) => {
                            const port =
                                container?.environment[u.port] ?? u.port;
                            const disabled = container.status !== "running";

                            // @ts-ignore
                            let url = new URL(window.apiURL);
                            url.port = port;
                            url.pathname = u.home ?? "";

                            return (
                                <a
                                    key={u.port}
                                    href={url.href}
                                    target="_blank"
                                    rel="noreferrer"
                                    className={classNames({
                                        [styles.navItem]: true,
                                        [styles.navItemDisabled]: disabled,
                                    })}
                                >
                                    <Horizontal>
                                        <MaterialIcon icon="public" />
                                        <Spacer />
                                        <MaterialIcon icon="open_in_new" />
                                    </Horizontal>
                                    <div>{url.href}</div>
                                </a>
                            );
                        })}
            </nav>
        </Content>
    );
}
