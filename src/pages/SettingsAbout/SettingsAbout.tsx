import { Fragment } from "react";
import { Title } from "../../components/Text/Text";
import Loading from "../../components/Loading/Loading";
import { useFetch } from "../../hooks/useFetch";
import { About } from "../../models/about";
import { api } from "../../backend/backend";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../components/KeyValueInfo/KeyValueInfo";

import styles from "./SettingsAbout.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(api.about.get);

    return (
        <Fragment>
            {loading && !error && (
                <Fragment>
                    <Title>Vertex</Title>
                    <Loading />
                </Fragment>
            )}
            {!loading && (
                <Vertical gap={30}>
                    <Vertical gap={20}>
                        <Title className={styles.title}>Vertex</Title>
                        <KeyValueGroup>
                            <KeyValueInfo
                                name="Version"
                                type="code"
                                symbol="tag"
                            >
                                {about?.version}
                            </KeyValueInfo>
                            <KeyValueInfo
                                name="Commit"
                                type="code"
                                symbol="commit"
                            >
                                {about?.commit}
                            </KeyValueInfo>
                            <KeyValueInfo
                                name="Release date"
                                type="code"
                                symbol="calendar_month"
                            >
                                {about?.date}
                            </KeyValueInfo>
                        </KeyValueGroup>
                    </Vertical>

                    <Vertical gap={20}>
                        <Title className={styles.title}>Platform</Title>
                        <KeyValueGroup>
                            <KeyValueInfo
                                name="OS"
                                type="code"
                                symbol="computer"
                            >
                                {about?.os}
                            </KeyValueInfo>
                            <KeyValueInfo
                                name="Architecture"
                                type="code"
                                symbol="memory"
                            >
                                {about?.arch}
                            </KeyValueInfo>
                        </KeyValueGroup>
                    </Vertical>
                </Vertical>
            )}
        </Fragment>
    );
}
