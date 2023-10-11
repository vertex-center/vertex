import { Fragment } from "react";
import { Title } from "../../../components/Text/Text";
import { api } from "../../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";

import styles from "./SettingsAbout.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useQuery } from "@tanstack/react-query";

export default function SettingsAbout() {
    const {
        data: about,
        isLoading,
        error,
    } = useQuery({
        queryKey: ["about"],
        queryFn: api.about,
    });

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <Vertical gap={20}>
                <Title className={styles.title}>Vertex</Title>
                <KeyValueGroup>
                    <KeyValueInfo
                        name="Version"
                        type="code"
                        icon="tag"
                        loading={isLoading}
                    >
                        {about?.version}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Commit"
                        type="code"
                        icon="commit"
                        loading={isLoading}
                    >
                        {about?.commit}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Release date"
                        type="code"
                        icon="calendar_month"
                        loading={isLoading}
                    >
                        {about?.date}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Compiled for"
                        type="code"
                        icon="memory"
                        loading={isLoading}
                    >
                        {about?.os}
                        {about?.arch && `/${about?.arch}`}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>
        </Fragment>
    );
}
