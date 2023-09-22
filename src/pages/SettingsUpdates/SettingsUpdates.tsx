import { useCallback, useEffect, useState } from "react";
import { Caption, Text, Title } from "../../components/Text/Text";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Button from "../../components/Button/Button";
import Spacer from "../../components/Spacer/Spacer";
import Symbol from "../../components/Symbol/Symbol";
import Popup from "../../components/Popup/Popup";
import Loading from "../../components/Loading/Loading";
import TimeAgo from "javascript-time-ago";
import en from "javascript-time-ago/locale/en";
import { Dependencies } from "../../models/update";
import { api } from "../../backend/backend";

import styles from "./SettingsUpdates.module.sass";
import Update, { Updates } from "../../components/Update/Update";
import { APIError } from "../../components/Error/Error";

TimeAgo.addDefaultLocale(en);

const timeAgo = new TimeAgo("en-US");

export default function SettingsUpdates() {
    const [dependencies, setDependencies] = useState<Dependencies>();
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState();
    const [updateError, setUpdateError] = useState();
    const [showMessage, setShowMessage] = useState<boolean>(false);

    const reload = useCallback((refresh?: boolean) => {
        setIsLoading(true);
        setError(undefined);

        api.dependencies
            .get(refresh)
            .then((res) => {
                setDependencies(res.data);
            })
            .catch(setError)
            .finally(() => setIsLoading(false));
    }, []);

    const updateService = (name: string) => {
        return api.dependencies
            .install([{ name }])
            .then(() => {
                if (name === "vertex") setShowMessage(true);
            })
            .catch(setUpdateError)
            .finally(reload);
    };

    useEffect(reload, []);

    const dismissPopup = () => {
        setShowMessage(false);
    };

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Updates</Title>
            {!isLoading && !error && (
                <Horizontal alignItems="center" gap={20}>
                    <Button onClick={() => reload(true)} rightSymbol="refresh">
                        Check for updates
                    </Button>
                    {dependencies?.last_updates_check && (
                        <Caption>
                            Last refresh:{" "}
                            {timeAgo.format(
                                new Date(dependencies?.last_updates_check),
                                "round"
                            )}
                        </Caption>
                    )}
                </Horizontal>
            )}
            {isLoading && <Loading />}
            {!error && !isLoading && dependencies?.items?.length === 0 && (
                <Horizontal alignItems="center">
                    <Symbol name="check" />
                    <Text>Everything is up-to-date.</Text>
                </Horizontal>
            )}
            <APIError error={updateError} />
            <Updates>
                {dependencies?.items?.map((dep) => (
                    <Update
                        key={dep?.id}
                        name={dep?.name}
                        version={dep?.version}
                        onUpdate={() => updateService(dep?.id)}
                        current_version={dep?.update?.current_version}
                        latest_version={dep?.update?.latest_version}
                        available={dep?.update !== undefined}
                    />
                ))}
            </Updates>
            <APIError error={error} />
            <Popup show={showMessage} onDismiss={dismissPopup}>
                <Text>
                    Updates are installed. You can now restart your Vertex
                    server.
                </Text>
                <Horizontal>
                    <Spacer />
                    <Button primary onClick={dismissPopup} rightSymbol="check">
                        OK
                    </Button>
                </Horizontal>
            </Popup>
        </Vertical>
    );
}
