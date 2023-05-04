import { Horizontal, Vertical } from "../Layouts/Layouts";

import styles from "./UptimeGraph.module.sass";

type Props = {
    title: string;
};

export default function UptimeGraph(props: Props) {
    const { title } = props;

    const count = 48;

    return (
        <Vertical gap={12} alignItems="flex-start">
            {/*<Horizontal>*/}
            {/*    <Text>{title}</Text>*/}
            {/*</Horizontal>*/}
            <Horizontal gap={6} className={styles.graph}>
                {[...new Array(count).keys()].map((i) => (
                    <div key={i} className={styles.bar} />
                ))}
            </Horizontal>
        </Vertical>
    );
}
