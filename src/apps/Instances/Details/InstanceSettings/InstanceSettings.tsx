import { Fragment, useEffect, useState } from "react";
import { Text, Title } from "../../../../components/Text/Text";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Button from "../../../../components/Button/Button";
import { useParams } from "react-router-dom";
import useInstance from "../../../../hooks/useInstance";
import Icon from "../../../../components/Icon/Icon";
import ToggleButton from "../../../../components/ToggleButton/ToggleButton";
import Input from "../../../../components/Input/Input";

import styles from "./InstanceSettings.module.sass";
import { api } from "../../../../backend/backend";
import { APIError } from "../../../../components/Error/APIError";
import Select, {
    SelectOption,
    SelectValue,
} from "../../../../components/Input/Select";
import VersionTag from "../../../../components/VersionTag/VersionTag";
import classNames from "classnames";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function InstanceSettings() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { instance, isLoading: isLoadingInstance } = useInstance(uuid);

    const [displayName, setDisplayName] = useState<string>();
    const [launchOnStartup, setLaunchOnStartup] = useState<boolean>();
    const [version, setVersion] = useState<string>();
    const [versions, setVersions] = useState<string[]>();
    const [versionsLoading, setVersionsLoading] = useState<boolean>(false);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [error, setError] = useState();

    useEffect(() => {
        if (!instance) return;
        setLaunchOnStartup(instance?.launch_on_startup ?? true);
        setDisplayName(instance?.display_name ?? instance?.service?.name);
        setVersion(instance?.version ?? "latest");
        reloadVersions();
    }, [instance]);

    const reloadVersions = (cache = true) => {
        setVersionsLoading(true);
        api.vxInstances
            .instance(instance.uuid)
            .versions.get(cache)
            .then((data) => {
                setVersions(data?.reverse());
            })
            .catch(setError)
            .finally(() => {
                setVersionsLoading(false);
            });
    };

    const mutationSave = useMutation({
        mutationFn: async () => {
            await api.vxInstances.instance(uuid).patch({
                launch_on_startup: launchOnStartup,
                display_name: displayName,
                version: version,
            });
        },
        onSuccess: () => {
            setSaved(true);
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["instances", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSave;

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
            <ProgressOverlay
                show={isLoadingInstance || versionsLoading || isUploading}
            />
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
                disabled={isLoadingInstance}
            />
            <div className={styles.versionSelect}>
                <Select
                    label="Version"
                    onChange={onVersionChange}
                    disabled={isLoadingInstance || versionsLoading}
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
                    rightIcon="refresh"
                    onClick={() => reloadVersions(false)}
                    disabled={isLoadingInstance || versionsLoading}
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
                    disabled={isLoadingInstance}
                />
            </Horizontal>
            <Button
                primary
                large
                onClick={async () => mutationSave.mutate()}
                rightIcon="save"
                loading={isUploading}
                disabled={saved || saved === undefined}
            >
                Save
            </Button>
            {!isUploading && saved && (
                <Horizontal
                    className={styles.saved}
                    alignItems="center"
                    gap={4}
                >
                    <Icon name="check" />
                    Saved!
                </Horizontal>
            )}
        </Fragment>
    );
}
