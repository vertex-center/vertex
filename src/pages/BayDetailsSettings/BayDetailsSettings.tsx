import { Fragment, useEffect, useState } from "react";
import { Text, Title } from "../../components/Text/Text";
import { Horizontal } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { patchInstance } from "../../backend/backend";
import { useParams } from "react-router-dom";
import useInstance from "../../hooks/useInstance";
import Progress from "../../components/Progress";
import Symbol from "../../components/Symbol/Symbol";
import ToggleButton from "../../components/ToggleButton/ToggleButton";

type Props = {};

export default function BayDetailsSettings(props: Props) {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const [autostart, setAutostart] = useState<boolean>();

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [uploading, setUploading] = useState<boolean>(false);
    const [error, setError] = useState<string>();

    useEffect(() => {
        console.log(instance);
        setAutostart(instance?.launch_on_startup ?? true);
    }, [instance]);

    const save = () => {
        setUploading(true);
        patchInstance(uuid, { launch_on_startup: autostart })
            .then(() => {
                setSaved(true);
            })
            .catch((error) => {
                setError(`${error.message}: ${error.response.data.message}`);
            })
            .finally(() => {
                setUploading(false);
            });
    };

    return (
        <Fragment>
            <Title>Settings</Title>
            <Horizontal alignItems="center">
                <Text>Launch on Startup</Text>
                <Spacer />
                <ToggleButton
                    value={autostart}
                    onChange={(v) => {
                        setAutostart(v);
                        setSaved(false);
                    }}
                />
            </Horizontal>
            <Button
                primary
                large
                onClick={save}
                rightSymbol="save"
                loading={uploading}
                disabled={saved || saved === undefined}
            >
                Save
            </Button>
            {uploading && <Progress infinite />}
            {!uploading && saved && (
                <Horizontal alignItems="center" gap={4}>
                    <Symbol name="check" />
                    Saved!
                </Horizontal>
            )}
            <Error error={error} />
        </Fragment>
    );
}
