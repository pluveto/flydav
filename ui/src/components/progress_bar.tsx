import styles from "./progress_bar.module.css";



interface ProgressBarProps {
    current: number,
    total: number
}
const ProgressBar = (props: ProgressBarProps) => {
    return (
        <div className={styles.barOuter}>
            <div className={styles.barInner} style={
                { width: (props.current / props.total * 100).toString() + "%" }
            }></div>
        </div>
    );

};
export default ProgressBar