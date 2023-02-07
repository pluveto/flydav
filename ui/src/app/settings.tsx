import styles from "./settings.module.css";
import sharedStyles from "../shared.module.css"
import { useEffect, useState } from "react";

interface SettingsProps {
    initialValue: SettingsObject
    onSave: (o: SettingsObject) => void
    onDiscard: () => void
}

export interface SettingsObject{
    url: string,
    username: string,
    password: string
}

const Settings = (props: SettingsProps) => {
    const [url, setUrl] = useState("");
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    
    useEffect(()=>{
        setUrl(props.initialValue.url)
        setUsername(props.initialValue.username)
        setPassword(props.initialValue.password)
    }, [props.initialValue])

    return (<div>
        <div className={styles.modal}>
            <div className={styles.modalInner}>
                <div className={styles.modalContent}>

                    <h3 className={styles.modalTitle}>Settings</h3>
                    <div className={styles.formEntry}>
                        <div className={styles.formEntryLabelOuter}>
                            <label className={styles.formEntryLabel}>
                                Server URL
                            </label>
                        </div>
                        <div className={styles.formEntryValueOuter}>
                            <input className={styles.formEntryInput} 
                            onInput={(e)=> setUrl((e.target as HTMLInputElement).value)}
                            type="text" value={url}></input>
                        </div>
                    </div>
                    <div className={styles.formEntry}>
                        <div className={styles.formEntryLabelOuter}>
                            <label className={styles.formEntryLabel}>
                                Username
                            </label>
                        </div>
                        <div className={styles.formEntryValueOuter}>
                            <input className={styles.formEntryInput}
                            onInput={(e)=> setUsername((e.target as HTMLInputElement).value)}
                            type="text" value={username}></input>
                        </div>
                    </div>
                    <div className={styles.formEntry}>
                        <div className={styles.formEntryLabelOuter}>
                            <label className={styles.formEntryLabel}>
                                Password
                            </label>
                        </div>
                        <div className={styles.formEntryValueOuter}>
                            <input className={styles.formEntryInput}
                            onInput={(e)=> setPassword((e.target as HTMLInputElement).value)}
                            type="password" value={password}></input>
                        </div>
                    </div>
                    <div className={styles.modalEnd}>
                        <button onClick={props.onDiscard} className={sharedStyles.buttonDefault}>Discard</button>
                        <button onClick={()=>{
                            props.onSave({url, username, password} as SettingsObject)
                        }} className={sharedStyles.buttonPrimary}>Save</button>
                    </div>

                </div>
            </div>
        </div>
        <div className={styles.mask}></div>
    </div>)
}

export default Settings