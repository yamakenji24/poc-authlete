import { Injectable } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { ConfigService } from '@nestjs/config';
import * as crypto from 'crypto';
import axios from 'axios';
import * as querystring from 'querystring';

@Injectable()
export class AuthService {
  private authDataStore = new Map<string, { codeVerifier: string; ticket: string }>();

  private readonly clientId: string;
  private readonly clientSecret: string;
  private readonly redirectUri: string;
  private readonly apiUrl: string;
  private readonly authleteToken: string;

  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService,
  ) {
    const baseURL = this.configService.get('AUTHLETE_BASE_URL') as string;
    const serviceId = this.configService.get('AUTHLETE_SERVICE_ID') as string;

    this.clientId = this.configService.get('AUTHLETE_CLIENT_ID') as string;
    this.clientSecret = this.configService.get(
      'AUTHLETE_CLIENT_SECRET',
    ) as string;
    this.redirectUri = this.configService.get(
      'AUTHLETE_REDIRECT_URI',
    ) as string;
    this.authleteToken = this.configService.get(
      'AUTHLETE_ACCESS_TOKEN',
    ) as string;
    this.apiUrl = `${baseURL}/${serviceId}`;
  }

  generateState(): string {
    return crypto.randomBytes(16).toString('hex');
  }

  generateCodeVerifier(): string {
    return crypto.randomBytes(32).toString('base64url');
  }

  generateCodeChallenge(codeVerifier: string): string {
    return crypto.createHash('sha256').update(codeVerifier).digest('base64url');
  }

  async getAuthorizationUrl(): Promise<string> {
    const codeVerifier = this.generateCodeVerifier();
    const codeChallenge = this.generateCodeChallenge(codeVerifier);
    const state = this.generateState();

    const params = querystring.stringify({
      response_type: 'code',
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: 'openid',
      code_challenge: codeChallenge,
      code_challenge_method: 'S256',
    });

    try {
      const response = await axios.post(
        `${this.apiUrl}/auth/authorization`,
        {
          parameters: params,
        },
        {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            Authorization: `Bearer ${this.authleteToken}`,
          },
        },
      );

      const ticket = response.data.ticket;
      this.authDataStore.set(state, { codeVerifier, ticket });
      return this.getLoginRedirectUri(state);
    } catch (error) {
      console.error(error);
      return '';
    }
  }

  async authenticateUser(email: string, password: string) {
    return true;
  }

  async getCallbackRedirectUri(ticket: string): Promise<string> {

    const response = await axios.post(`${this.apiUrl}/auth/authorization/issue`, {
      ticket: ticket,
      subject: 'yamakenji',
    }, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        Authorization: `Bearer ${this.authleteToken}`,
      },
    });

    return response.data.responseContent;
  }

  async exchangeCodeForTokens(code: string, codeVerifier: string) {
    const params = querystring.stringify({
      grant_type: 'authorization_code',
      code,
      redirect_uri: this.redirectUri,
      code_verifier: codeVerifier,
    });

    try {
      const response = await axios.post(`${this.apiUrl}/auth/token`, {
        parameters: params,
        clientId: this.clientId,
        clientSecret: this.clientSecret,
      }, {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.authleteToken}`,
        },
      });

      return response.data;
    } catch (error) {
      console.error(error);
    }
  }

  getLoginRedirectUri(state: string) {
    return `http://localhost:3000/login/identify?state=${state}`;
  }
  getAuthData(state: string) {
    return this.authDataStore.get(state);
  }
}
