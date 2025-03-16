import { Controller, Get, Post, Query, Req, Res } from '@nestjs/common';
import { AuthService } from './auth.service';
import { Request, Response } from 'express';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Get('authorize')
  async authorize(@Res() res: Response) {
    const authorizationUrl = await this.authService.getAuthorizationUrl();

    console.log('authorizationUrl:', authorizationUrl);

    return res.redirect(authorizationUrl);
  }

  @Post('/login')
  async login(@Req() req: Request, @Res() res: Response) {
    const state = req.body.state;
    const email = req.body.email;
    const password = req.body.password;

    try {
      this.authService.authenticateUser(email, password);
      const authData = this.authService.getAuthData(state);
      if (!authData) {
        throw new Error('AuthData not found');
      }
      const { ticket } = authData;
      const redirectUri = new URL(await this.authService.getCallbackRedirectUri(ticket));
      redirectUri.searchParams.append('state', state);
      return res.redirect(redirectUri.toString());
    } catch (error) {
      console.error(error);

      return res.send('Error');
    }
  }

  @Get('/callback')
  async callback(@Query('state') state: string, @Query('code') code: string, @Res() res: Response) {
    if (!state || !code) {
      return res.send('Invalid request');
    }

    const authData = this.authService.getAuthData(state);
    if (!authData) {
      return res.send('AuthData not found');
    }

    const { codeVerifier } = authData;

    const tokens = await this.authService.exchangeCodeForTokens(code, codeVerifier);

    console.log('tokens:', tokens);

    return res.send(tokens);
  }
}
