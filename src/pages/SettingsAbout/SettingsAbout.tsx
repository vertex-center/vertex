import { Fragment } from "react";
import { Title } from "../../components/Text/Text";
import { getAbout } from "../../backend/backend";
import Loading from "../../components/Loading/Loading";
import Symbol from "../../components/Symbol/Symbol";

import styles from "./SettingsAbout.module.sass";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import { useFetch } from "../../hooks/useFetch";
import { About } from "../../models/about";

type ItemProps = {
    symbol: string;
    label: string;
    value?: string;
};

function Item(props: ItemProps) {
    const { symbol, label } = props;
    let { value } = props;

    if (!value) value = "N/A";

    return (
        <Horizontal gap={12} alignItems="center">
            <div className={styles.item}>
                <Symbol name={symbol} />
            </div>
            <div className={styles.item}>{label}:</div>
            <div className={styles.item}>{value}</div>
        </Horizontal>
    );
}

export default function SettingsAbout() {
    const { data: about, loading, error } = useFetch<About>(getAbout);

    return (
        <Fragment>
            {/*{error && <Error error={error} />}*/}
            {loading && !error && (
                <Fragment>
                    <Title>Vertex</Title>
                    <Loading />
                </Fragment>
            )}
            {!loading && (
                <Fragment>
                    <Title>Vertex</Title>
                    <Vertical gap={4} style={{ marginBottom: 16 }}>
                        <Item
                            symbol="tag"
                            label="Version"
                            value={about?.version}
                        />
                        <Item
                            symbol="commit"
                            label="Commit"
                            value={about?.commit}
                        />
                        <Item
                            symbol="calendar_month"
                            label="Release date"
                            value={about?.date}
                        />
                    </Vertical>
                    <Title>Platform</Title>
                    <Vertical gap={4}>
                        <Item symbol="computer" label="OS" value={about?.os} />
                        <Item
                            symbol="memory"
                            label="Arch"
                            value={about?.arch}
                        />
                    </Vertical>
                </Fragment>
            )}
        </Fragment>
    );
}
