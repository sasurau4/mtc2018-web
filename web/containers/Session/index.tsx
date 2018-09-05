import * as React from 'react';
import Head from 'next/head';
import styled from 'styled-components';
import Router, { withRouter } from 'next/router';
import Default from '../../layout/Default';
import { Content } from '../../types';
import ContentCard from './ContentCard';
import Header from './Header';
import { Button, Section } from '../../components';
import { withI18next } from '../../lib/with-i18next';
import { I18n } from 'react-i18next';

/* tslint:disable-next-line:no-var-requires */
const contentsData = require('../../static/json/contents.json');

interface Props {
  router: any; // TODO
}

class Session extends React.Component<Props> {
  private content!: Content;

  public componentWillMount() {
    // コンテンツ内容を取得
    const contentId = Number(this.props.router!.query!.id);
    this.content = contentsData.sessions.find((content: Content) => {
      return content.id === contentId;
    })!;
  }

  public render() {
    return (
      <Default>
        <I18n>
          {(_, { i18n }) => {
            const isJa = i18n.language === 'ja-JP';
            return (
              <>
                <Head>
                  <title>
                    Mercari Tech Conf 2018 -{' '}
                    {isJa ? this.content.titleJa : this.content.title}
                  </title>
                </Head>
                <Header />
                <Body>
                  <Section title="SESSION">
                    <ContentCard content={this.content} isJa={isJa} />
                    <BackButton
                      type="primary"
                      size="large"
                      onClick={this.onClickBackButton}
                    >
                      BACK
                    </BackButton>
                  </Section>
                </Body>
              </>
            );
          }}
        </I18n>
      </Default>
    );
  }

  private onClickBackButton = () => {
    Router.back();
  };
}

const Body = styled.div`
  width: 100%;
  padding: 32px 64px 64px;
  box-sizing: border-box;

  > * {
    margin-bottom: 160px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  @media screen and (max-width: 767px) {
    padding: 32px 8px;

    > * {
      margin-bottom: 80px;

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
`;

const BackButton = styled(Button)`
  width: 200px;
  margin-top: 60px;
`;

export default withI18next()(withRouter(Session));