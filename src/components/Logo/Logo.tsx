import styles from "./Logo.module.sass";

type Props = {
    iconOnly?: boolean;
};

export default function Logo({ iconOnly }: Readonly<Props>) {
    return (
        <div className={styles.logo}>
            <img width={30} src="/images/logo.png" alt="Logo" />
            {!iconOnly && <h1>Vertex</h1>}
        </div>
    );
}
