import { Fragment, useEffect, useState } from "react";
import { Text, Title } from "../../../../components/Text/Text";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Button from "../../../../components/Button/Button";
import { useParams } from "react-router-dom";
import useInstance from "../../../../hooks/useInstance";
import Progress from "../../../../components/Progress";
import Symbol from "../../../../components/Symbol/Symbol";
import ToggleButton from "../../../../components/ToggleButton/ToggleButton";
import Input from "../../../../components/Input/Input";

import styles from "./InstanceSettings.module.sass";
import { api } from "../../../../backend/backend";
import { APIError } from "../../../../components/Error/Error";

type Props = {};

export default function InstanceSettings(props: Props) {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const [displayName, setDisplayName] = useState<string>();
    const [launchOnStartup, setLaunchOnStartup] = useState<boolean>();

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [uploading, setUploading] = useState<boolean>(false);
    const [error, setError] = useState();

    useEffect(() => {
        if (!instance) return;
        setLaunchOnStartup(instance?.launch_on_startup ?? true);
        setDisplayName(instance?.display_name ?? instance?.service?.name);
    }, [instance]);

    const save = () => {
        setUploading(true);
        api.instance
            .patch(uuid, {
                launch_on_startup: launchOnStartup,
                display_name: displayName,
            })
            .then(() => {
                setSaved(true);
            })
            .catch(setError)
            .finally(() => {
                setUploading(false);
            });
    };

    return (
        <Fragment>
            <Title className={styles.title}>Settings</Title>
            <APIError error={error} />
            <Input
                label="Instance name"
                description="The custom name of your choice for this service"
                value={displayName}
                onChange={(e: any) => {
                    setDisplayName(e.target.value);
                    setSaved(false);
                }}
                disabled={displayName === undefined}
            />
            <Horizontal className={styles.toggle} alignItems="center">
                <Text>Launch on Startup</Text>
                <Spacer />
                <ToggleButton
                    value={launchOnStartup}
                    onChange={(v) => {
                        setLaunchOnStartup(v);
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
                <Horizontal
                    className={styles.saved}
                    alignItems="center"
                    gap={4}
                >
                    <Symbol name="check" />
                    Saved!
                </Horizontal>
            )}
        </Fragment>
    );
}
