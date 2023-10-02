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
import Select, {
    SelectOption,
    SelectValue,
} from "../../../../components/Input/Select";
import VersionTag from "../../../../components/VersionTag/VersionTag";
import classNames from "classnames";

type Props = {};

export default function InstanceSettings(props: Props) {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const [displayName, setDisplayName] = useState<string>();
    const [launchOnStartup, setLaunchOnStartup] = useState<boolean>();
    const [version, setVersion] = useState<string>();
    const [versions, setVersions] = useState<string[]>();

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [uploading, setUploading] = useState<boolean>(false);
    const [error, setError] = useState();

    useEffect(() => {
        if (!instance) return;
        setLaunchOnStartup(instance?.launch_on_startup ?? true);
        setDisplayName(instance?.display_name ?? instance?.service?.name);
        setVersion(instance?.version ?? "latest");
        reloadVersions();
    }, [instance]);

    const reloadVersions = (cache = true) => {
        api.instance.versions
            .get(instance.uuid)
            .then((res) => {
                setVersions(res.data?.reverse());
            })
            .catch(setError);
    };

    const save = () => {
        setUploading(true);
        api.instance
            .patch(uuid, {
                launch_on_startup: launchOnStartup,
                display_name: displayName,
                version: version,
            })
            .then(() => {
                setSaved(true);
            })
            .catch(setError)
            .finally(() => {
                setUploading(false);
            });
    };

    const onVersionChange = (v: any) => {
        setVersion(v);
        setSaved(false);
    };

    const versionValue = (
        <SelectValue
            className={classNames({
                [styles.versionValue]: version !== "latest",
            })}
        >
            {version === "latest" ? (
                "Always pull latest version"
            ) : (
                <VersionTag>{version}</VersionTag>
            )}
        </SelectValue>
    );

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
            <div className={styles.versionSelect}>
                <Select
                    label="Version"
                    onChange={onVersionChange}
                    // @ts-ignore
                    value={versionValue}
                >
                    {versions?.includes("latest") && (
                        <SelectOption value="latest">
                            Always pull latest version
                        </SelectOption>
                    )}
                    {versions?.map((v) => {
                        if (v === "latest") {
                            return null;
                        }
                        return (
                            <SelectOption
                                key={v}
                                value={v}
                                className={styles.versionOption}
                            >
                                <VersionTag>{v}</VersionTag>
                            </SelectOption>
                        );
                    })}
                </Select>
                <Button
                    rightSymbol="refresh"
                    onClick={() => reloadVersions(true)}
                >
                    Refresh
                </Button>
            </div>
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
