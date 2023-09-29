import { Fragment } from "react";
import { Title } from "../../../components/Text/Text";
import Loading from "../../../components/Loading/Loading";
import { useFetch } from "../../../hooks/useFetch";
import { About } from "../../../models/about";
import { api } from "../../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";

import styles from "./SettingsAbout.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/Error";

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(api.about.get);

    return (
        <Fragment>
            {(loading || error) && (
                <Vertical>
                    <Title className={styles.title}>Vertex</Title>
                    {loading && <Loading />}
                </Vertical>
            )}
            <APIError error={error} />
            {!loading && !error && (
                <Vertical gap={20}>
                    <Title className={styles.title}>Vertex</Title>
                    <KeyValueGroup>
                        <KeyValueInfo name="Version" type="code" symbol="tag">
                            {about?.version}
                        </KeyValueInfo>
                        <KeyValueInfo name="Commit" type="code" symbol="commit">
                            {about?.commit}
                        </KeyValueInfo>
                        <KeyValueInfo
                            name="Release date"
                            type="code"
                            symbol="calendar_month"
                        >
                            {about?.date}
                        </KeyValueInfo>
                        <KeyValueInfo
                            name="Compiled for"
                            type="code"
                            symbol="memory"
                        >
                            {about?.os}/{about?.arch}
                        </KeyValueInfo>
                    </KeyValueGroup>
                </Vertical>
            )}
        </Fragment>
    );
}
