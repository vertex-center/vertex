import { useEffect, useState } from "react";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import {
    Button,
    MaterialIcon,
    Paragraph,
    SelectField,
    SelectOption,
    TextField,
    Title,
} from "@vertex-center/components";
import { useParams } from "react-router-dom";
import useContainer from "../../hooks/useContainer";
import ToggleButton from "../../../../components/ToggleButton/ToggleButton";
import styles from "./ContainerSettings.module.sass";
import { APIError } from "../../../../components/Error/APIError";
import VersionTag from "../../../../components/VersionTag/VersionTag";
import classNames from "classnames";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";

export default function ContainerSettings() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { container, isLoading: isLoadingContainer } = useContainer(uuid);

    const [name, setName] = useState<string>();
    const [launchOnStartup, setLaunchOnStartup] = useState<boolean>();
    const [imageTag, setImageTag] = useState<string>();
    const [imageTags, setImageTags] = useState<string[]>();
    const [imageTagsLoading, setImageTagsLoading] = useState<boolean>(false);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [error, setError] = useState();

    useEffect(() => {
        if (!container) return;
        setLaunchOnStartup(container.launch_on_startup);
        setName(container.name);
        setImageTag(container?.image_tag ?? "latest");
        reloadVersions();
    }, [container]);

    const reloadVersions = (cache = true) => {
        setImageTagsLoading(true);
        API.getVersions(container.id, cache)
            .then((data) => {
                setImageTags(data?.reverse());
            })
            .catch(setError)
            .finally(() => {
                setImageTagsLoading(false);
            });
    };

    const mutationSave = useMutation({
        mutationFn: async () => {
            await API.patchContainer(uuid, {
                launch_on_startup: launchOnStartup,
                name,
                imageTag,
            });
        },
        onSuccess: () => {
            setSaved(true);
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSave;

    const onVersionChange = (v: any) => {
        setImageTag(v);
        setSaved(false);
    };

    const versionValue = (
        <div
            className={classNames({
                [styles.versionValue]: imageTag !== "latest",
            })}
        >
            {imageTag === "latest" ? (
                "Always pull latest"
            ) : (
                <VersionTag>{imageTag}</VersionTag>
            )}
        </div>
    );

    return (
        <Content>
            <Title variant="h2">Settings</Title>
            <ProgressOverlay
                show={isLoadingContainer || imageTagsLoading || isUploading}
            />
            <APIError error={error} />
            <Horizontal alignItems="center">
                <Paragraph>Launch on Startup</Paragraph>
                <Spacer />
                <ToggleButton
                    value={launchOnStartup}
                    onChange={(v) => {
                        setLaunchOnStartup(v);
                        setSaved(false);
                    }}
                    disabled={isLoadingContainer}
                />
            </Horizontal>
            <TextField
                id="container-name"
                label="Container name"
                description="The custom name of your choice for this service"
                value={name}
                onChange={(e: any) => {
                    setName(e.target.value);
                    setSaved(false);
                }}
                disabled={isLoadingContainer}
            />
            <div className={styles.versionSelect}>
                <SelectField
                    id="container-version"
                    label="Version"
                    onChange={onVersionChange}
                    disabled={isLoadingContainer || imageTagsLoading}
                    // @ts-ignore
                    value={versionValue}
                >
                    {imageTags?.includes("latest") && (
                        <SelectOption value="latest">
                            Always pull latest
                        </SelectOption>
                    )}
                    {imageTags?.map((v) => {
                        if (v === "latest") {
                            return null;
                        }
                        return (
                            <SelectOption key={v} value={v}>
                                <VersionTag>{v}</VersionTag>
                            </SelectOption>
                        );
                    })}
                </SelectField>
                <Button
                    rightIcon={<MaterialIcon icon="refresh" />}
                    onClick={() => reloadVersions(false)}
                    disabled={isLoadingContainer || imageTagsLoading}
                >
                    Refresh
                </Button>
            </div>
            <Horizontal justifyContent="flex-end">
                {!isUploading && saved && (
                    <Horizontal
                        className={styles.saved}
                        alignItems="center"
                        gap={4}
                    >
                        <MaterialIcon icon="check" />
                        Saved!
                    </Horizontal>
                )}
                <Button
                    variant="colored"
                    onClick={async () => mutationSave.mutate()}
                    rightIcon={<MaterialIcon icon="save" />}
                    disabled={isUploading || saved || saved === undefined}
                >
                    Save
                </Button>
            </Horizontal>
        </Content>
    );
}
