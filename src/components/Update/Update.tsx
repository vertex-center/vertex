import { PropsWithChildren, useState } from "react";

import styles from "./Update.module.sass";
import Button from "../Button/Button";
import Symbol from "../Symbol/Symbol";
import Progress from "../Progress";

export function Updates(props: PropsWithChildren) {
    return <div {...props} />;
}

type Props = {
    name: string;
    version?: string;
    available?: boolean;

    current_version?: string;
    latest_version?: string;

    onUpdate: () => Promise<void>;
};

export default function Update(props: Props) {
    const { name, available, version, latest_version } = props;

    const [isLoading, setIsLoading] = useState(false);

    const onUpdate = () => {
        setIsLoading(true);
        props.onUpdate().finally(() => setIsLoading(false));
    };

    return (
        <div className={styles.update}>
            <div className={styles.info}>
                <div className={styles.name}>{name}</div>
                {version && (
                    <code className={styles.version}>
                        {version}
                        {latest_version && " -> " + latest_version}
                    </code>
                )}
            </div>
            {available ? (
                isLoading ? (
                    <Progress infinite />
                ) : (
                    <Button rightSymbol="download" onClick={onUpdate}>
                        Update
                    </Button>
                )
            ) : (
                <div className={styles.status}>
                    <Symbol name="check" />
                    Up-to-date
                </div>
            )}
        </div>
    );
}
