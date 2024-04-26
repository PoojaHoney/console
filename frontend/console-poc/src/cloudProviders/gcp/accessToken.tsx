import axios from 'axios';
const getAccessToken = async(jwtToken: string): Promise<string> => {
    debugger
    try{
        const response = await axios.post(process.env.REACT_APP_GCP_TOKEN_URL!, {
            grant_type: "urn:ietf:params:oauth:grant-type:jwt-bearer",
            assertion: jwtToken
        });
        return response.data.access_token;
    } catch(error: any) {
        throw new Error('Error exchanging JWT fot access token: ' + error.response.data.error_description);
    }
};

const generateJWT = async(): Promise<string> => {
    debugger
    const jwt = require('jsonwebtoken');

    const jwtPayload = {
        iss: process.env.REACT_APP_GCP_SRV_ACC_CLIENT_EMAIL,
        sub: process.env.REACT_APP_GCP_SRV_ACC_CLIENT_EMAIL,
        aud: process.env.REACT_APP_GCP_TOKEN_URL,
        exp: Math.floor(Date.now()/1000)+3600,
    }
    const private_key = process.env.REACT_APP_GCP_SRV_ACC_PRIVATE_KEY;
    const jwtToken = jwt.sign(jwtPayload, private_key, { algorithm: 'RS256' });
    return jwtToken
}

export {getAccessToken, generateJWT}