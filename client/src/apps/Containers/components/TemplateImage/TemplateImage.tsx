import ServiceLogo from "../../../../components/ServiceLogo/ServiceLogo";
import styles from "./TemplateImage.module.sass";

type Props = {
    icon: string;
    color: string;
};

export default function TemplateImage(props: Readonly<Props>) {
    const { icon, color } = props;

    const backgroundColor = `${color ?? "#000000"}08`;

    return (
        <div className={styles.logo} style={{ backgroundColor }}>
            <ServiceLogo icon={icon} color={color} />
        </div>
    );
}
