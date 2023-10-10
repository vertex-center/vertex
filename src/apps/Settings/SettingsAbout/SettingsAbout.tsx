import { Fragment } from "react";
import { Title } from "../../../components/Text/Text";
import { useFetch } from "../../../hooks/useFetch";
import { About } from "../../../models/about";
import { api } from "../../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";

import styles from "./SettingsAbout.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(api.about);

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <APIError error={error} />
            <Vertical gap={20}>
                <Title className={styles.title}>Vertex</Title>
                <KeyValueGroup>
                    <KeyValueInfo
                        name="Version"
                        type="code"
                        icon="tag"
                        loading={loading}
                    >
                        {about?.version}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Commit"
                        type="code"
                        icon="commit"
                        loading={loading}
                    >
                        {about?.commit}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Release date"
                        type="code"
                        icon="calendar_month"
                        loading={loading}
                    >
                        {about?.date}
                    </KeyValueInfo>
                    <KeyValueInfo
                        name="Compiled for"
                        type="code"
                        icon="memory"
                        loading={loading}
                    >
                        {about?.os}
                        {about?.arch && `/${about?.arch}`}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>
        </Fragment>
    );
}
