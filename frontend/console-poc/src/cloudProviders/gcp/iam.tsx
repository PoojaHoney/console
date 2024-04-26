import React, { useState } from 'react';
import axios from 'axios';
import {generateJWT, getAccessToken} from './accessToken';
import {Button, Input, message} from 'antd';

const CreateServiceAccount: React.FC = () => {
    const [projectID, setProjectID] = useState<string>(process.env.REACT_APP_GCP_PROJECT_ID || '');
    const [serviceAccountName, setServiceAccountName] = React.useState('');
    const [messageText, setMessageText] = React.useState('');
    debugger
    const handleCreateServiceAccount = async () => {
        debugger
        try {
            const jwtToken = await generateJWT();
            const accessToken = await getAccessToken(jwtToken);

            axios.post(`https://iam.googleapis.com/v1/projects/${projectID}/serviceAccounts`, {
                accountId: serviceAccountName,
                serviceAccount: {
                    displayName: serviceAccountName
                }
            }, {
                headers: {
                    Authorization: `Bearer ${accessToken}`
                }
            }).then((response) => {
                console.log("service account created: ", response.data);
                setMessageText('Service account created successfully');
                message.success('Service account created successfully');
            }).catch((error) => {
                console.log("error creating service account: ", error);
                setMessageText('Error creating service account');
                message.error('Error creating service account');
            })
        } catch(error) {
            console.log("error creating service account: ", error);
            setMessageText('Error creating service account');
            message.error('Error creating service account');
        }
    }

    return (
        <>
            <Input placeholder='service account name' value={serviceAccountName} onChange={(e) => setServiceAccountName(e.target.value)} />
            <Button onClick={handleCreateServiceAccount}>Create Service Account</Button>
            <>{messageText}</>
        </>
    )
}

export default CreateServiceAccount;
