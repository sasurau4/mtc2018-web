import * as React from 'react';
import Head from 'next/head';
import styled from 'styled-components';
import Router, { withRouter, WithRouterProps } from 'next/router';
import Default from '../../../layout/Default';
import ContentCard, {
  CONTENT_CARD_FRAGMENT
} from '../../../containers/Session/ContentCard';
import Header from '../../../containers/Session/Header';
import { Button, Section } from '../../../components';
import { withI18next } from '../../../lib/with-i18next';
import { NamespacesConsumer } from 'react-i18next';
import { isJapan, generateUUID } from '../../../utils';

import gql from 'graphql-tag';
import { Query, Mutation } from 'react-apollo';
import { FetchResult } from 'apollo-link';
import { DataProxy } from 'apollo-cache';
import {
  Session as SessionQuery,
  SessionVariables
} from '../../../graphql/generated/Session';
import { SessionFragment } from '../../../graphql/generated/SessionFragment';

import { LIKE_SESSION_MUTATION } from '../../../graphql/mutation';
import {
  LikeSession,
  LikeSessionVariables
} from '../../../graphql/generated/LikeSession';

export const SESSION_QUERY = gql`
  query Session($sessionId: Int!) {
    session(sessionId: $sessionId) {
      ...SessionFragment
      ...ContentCardFragment
    }
  }

  fragment SessionFragment on Session {
    title
    titleJa
  }

  ${CONTENT_CARD_FRAGMENT}
`;

class SessionQueryComponent extends Query<SessionQuery, SessionVariables> {}

class CreateLikeMutationComponent extends Mutation<
  LikeSession,
  LikeSessionVariables
> {}

class Session extends React.Component<WithRouterProps> {
  public render() {
    const sessionId = parseInt(this.props.router!.query!.id as string, 10);
    return (
      <Default>
        <SessionQueryComponent
          query={SESSION_QUERY}
          variables={{ sessionId }}
          pollInterval={1000}
        >
          {({ data, error, loading }) => {
            if (error || loading || !data || !data.session) {
              return null;
            }

            return (
              <CreateLikeMutationComponent
                mutation={LIKE_SESSION_MUTATION}
                update={this.update(sessionId)}
                variables={{ sessionId, randomId: generateUUID() }}
              >
                {createLike => (
                  <NamespacesConsumer ns={['common']}>
                    {(_, { i18n }) => {
                      const isJa = isJapan(i18n.language);
                      const session: SessionFragment = data.session!;
                      return (
                        <>
                          <Head>
                            <title>
                              Mercari Tech Conf 2018 -{' '}
                              {isJa ? session.titleJa : session.title}
                            </title>
                          </Head>
                          <StyledHeader />
                          <Body>
                            <Section title="SESSION">
                              <ContentCard
                                session={data.session!}
                                isJa={isJa}
                                onClickLikeButton={createLike}
                              />
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
                  </NamespacesConsumer>
                )}
              </CreateLikeMutationComponent>
            );
          }}
        </SessionQueryComponent>
      </Default>
    );
  }

  private onClickBackButton = () => {
    Router.back();
  };

  private update = (sessionId: number) => (
    cache: DataProxy,
    mutationResult: FetchResult<LikeSession>
  ) => {
    // TODO: How to define types?
    // This is dirty and not good solution.
    const cachedData = cache.readQuery({
      query: SESSION_QUERY,
      variables: { sessionId }
    })! as { session: SessionQuery };

    const updatedSession = {
      ...cachedData.session,
      liked: mutationResult.data!.createLike!.like.session.liked
    };
    cache.writeQuery({
      query: SESSION_QUERY,
      variables: { sessionId },
      data: { session: updatedSession }
    });
  };
}

const StyledHeader = styled(Header)`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
`;

const Body = styled.div`
  width: 100%;
  padding: 96px 64px 64px;
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
