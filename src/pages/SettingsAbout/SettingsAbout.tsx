import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { getAbout } from "../../backend/backend";
import Loading from "../../components/Loading/Loading";

type Props = {};

export default function SettingsAbout(props: Props) {
    const [version, setVersion] = useState<string>();
    const [commit, setCommit] = useState<string>();
    const [date, setDate] = useState<string>();

    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        setLoading(true);
        getAbout().then((about) => {
            setVersion(about.version);
            setCommit(about.commit);
            setDate(about.date);
            setLoading(false);
        });
    }, []);

    return (
        <Fragment>
            <Title>About</Title>
            {loading && <Loading />}
            {!loading && (
                <Fragment>
                    <div>Version: {version}</div>
                    <div>Commit: {commit}</div>
                    <div>Release date: {date}</div>
                </Fragment>
            )}
        </Fragment>
    );
}
