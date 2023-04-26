import { Fragment, useCallback, useEffect, useState } from "react";
import { Text, Title } from "../../components/Text/Text";
import { Horizontal } from "../../components/Layouts/Layouts";
import Button from "../../components/Button/Button";
import Spacer from "../../components/Spacer/Spacer";
import Symbol from "../../components/Symbol/Symbol";
import { executeUpdates, getUpdates } from "../../backend/backend";
import { Error } from "../../components/Error/Error";
import Progress from "../../components/Progress";
import Popup from "../../components/Popup/Popup";
import Loading from "../../components/Loading/Loading";

type Props = {
    name: string;
    update: () => void;
    current_version: string;
    latest_version: string;
    upToDate: boolean;
};

function Update(props: Props) {
    const { name, update, current_version, latest_version, upToDate } = props;

    const [isLoading, setIsLoading] = useState(false);

    const onUpdate = () => {
        setIsLoading(true);
        update();
    };

    return (
        <Horizontal gap={24} alignItems="center">
            <Text>{name}</Text>
            <Spacer />
            <code>
                {current_version} {!upToDate && "->"}{" "}
                {!upToDate && latest_version}
            </code>
            {!upToDate && !isLoading && (
                <Button onClick={onUpdate} rightSymbol="download">
                    Update
                </Button>
            )}
            {!upToDate && isLoading && <Progress infinite />}
            {upToDate && (
                <Horizontal
                    gap={6}
                    alignItems="center"
                    style={{ color: "var(--green)" }}
                >
                    <Symbol name="check" />
                    Up to date
                </Horizontal>
            )}
        </Horizontal>
    );
}

export default function SettingsUpdates() {
    const [updates, setUpdates] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState();
    const [showMessage, setShowMessage] = useState<boolean>(false);

    const reload = useCallback(() => {
        setIsLoading(true);
        setUpdates(null);
        getUpdates()
            .then((updates) => {
                setUpdates(updates);
            })
            .catch((err) => {
                setError(err?.response?.data?.message ?? err?.message);
            })
            .finally(() => setIsLoading(false));
    }, []);

    const updateService = (name: string) => {
        executeUpdates([{ name }])
            .then(() => setShowMessage(true))
            .catch((err) => {
                setError(err?.response?.data?.message ?? err?.message);
            })
            .finally(reload);
    };

    useEffect(reload, []);

    const dismissPopup = () => {
        setShowMessage(false);
    };

    return (
        <Fragment>
            <Title>Updates</Title>
            {!error && !isLoading && updates === null && (
                <Horizontal alignItems="center">
                    <Symbol name="check" />
                    <Text>Everything is up-to-date.</Text>
                </Horizontal>
            )}
            {isLoading && <Loading />}
            {updates?.map((update) => (
                <Update
                    key={update.name}
                    name={update.name}
                    latest_version={update.latest_version}
                    current_version={update.current_version}
                    upToDate={update.up_to_date}
                    update={() => updateService(update.id)}
                />
            ))}
            {error && <Error error={error} />}
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
        </Fragment>
    );
}
