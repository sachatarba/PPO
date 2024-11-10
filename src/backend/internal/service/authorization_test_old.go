package service

// import (
// 	"context"

// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/sachatarba/course-db/internal/entity"
// 	"github.com/sachatarba/course-db/internal/service/mocks"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// var testClientID = func() uuid.UUID {
// 	b, _ := uuid.ParseBytes([]byte("aaaaaaaaaaaaaaaa"))
// 	return b
// }()

// var testSessionID = func() uuid.UUID {
// 	b, _ := uuid.ParseBytes([]byte("bbbbbbbbbbbbbbbb"))
// 	return b
// }()

// func TestAuthorizationService_Authorize(t *testing.T) {
// 	type sessionRepoMockParams struct {
// 		ctx     any
// 		session any
// 		err     any
// 	}
// 	type clientServiceMockParams struct {
// 		ctx    any
// 		login  any
// 		client any
// 		err    any
// 	}
// 	type args struct {
// 		ctx      context.Context
// 		login    string
// 		password string
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		sessionMock sessionRepoMockParams
// 		clientMock  clientServiceMockParams
// 		want        entity.Session
// 		wantErr     error
// 	}{
// 		{
// 			name: "Authorization ok",
// 			args: args{
// 				ctx:      context.Background(),
// 				login:    "sus",
// 				password: "123",
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     nil,
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx:   mock.Anything,
// 				login: "sus",
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 				err: nil,
// 			},
// 			want: entity.Session{
// 				ClientID: testClientID,
// 			},
// 			wantErr: nil,
// 		},

// 		{
// 			name: "Authorization failed: no user",
// 			args: args{
// 				ctx:      context.Background(),
// 				login:    "sus",
// 				password: "123",
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     nil,
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx:    mock.Anything,
// 				login:  "sus",
// 				client: entity.Client{},
// 				err:    ErrNoSuchClient,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrNoSuchClient,
// 		},

// 		{
// 			name: "Authorization failed: wrong password",
// 			args: args{
// 				ctx:      context.Background(),
// 				login:    "sus",
// 				password: "123",
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     nil,
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx:   mock.Anything,
// 				login: "sus",
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "1234",
// 				},
// 				err: ErrWrongPassword,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrWrongPassword,
// 		},

// 		{
// 			name: "Authorization failed: sessionRepo error",
// 			args: args{
// 				ctx:      context.Background(),
// 				login:    "sus",
// 				password: "123",
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     ErrInternalSessionRepo,
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx:   mock.Anything,
// 				login: mock.Anything,
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 				err: nil,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrInternalSessionRepo,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sessionRepo := &mocks.ISessionRepository{}
// 			clientService := &mocks.IClientService{}

// 			sessionRepo.
// 				On("CreateNewSession", tt.sessionMock.ctx, tt.sessionMock.session).
// 				Return(tt.sessionMock.err)
// 			clientService.
// 				On("GetClientByLogin", tt.clientMock.ctx, tt.clientMock.login).
// 				Return(tt.clientMock.client, tt.clientMock.err)

// 			a := &AuthorizationService{
// 				sessionRepo:   sessionRepo,
// 				clientService: clientService,
// 			}

// 			got, err := a.Authorize(tt.args.ctx, tt.args.login, tt.args.password)

// 			assert.ErrorIsf(t, err, tt.wantErr, "errors misplace")

// 			assert.Equal(t, tt.want.ClientID, got.ClientID, "wrong client id")
// 		})
// 	}
// }

// func TestAuthorizationService_Register(t *testing.T) {
// 	type clientServiceMockParams struct {
// 		ctx    any
// 		client any
// 		err    any
// 	}
// 	type sessionRepoMockParams struct {
// 		ctx     any
// 		session any
// 		err     any
// 	}
// 	type args struct {
// 		ctx    context.Context
// 		client entity.Client
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		sessionMock sessionRepoMockParams
// 		clientMock  clientServiceMockParams
// 		want        entity.Session
// 		wantErr     error
// 	}{
// 		{
// 			name: "Registration ok",
// 			args: args{
// 				ctx: context.Background(),
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx: mock.Anything,
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 				err: nil,
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     nil,
// 			},
// 			want: entity.Session{
// 				ClientID: testClientID,
// 			},
// 			wantErr: nil,
// 		},

// 		{
// 			name: "Registration fail: clientRepo error",
// 			args: args{
// 				ctx: context.Background(),
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx: mock.Anything,
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 				err: ErrInternalDB,
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     nil,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrInternalDB,
// 		},

// 		{
// 			name: "Registration fail: sessionRepo error",
// 			args: args{
// 				ctx: context.Background(),
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 			},
// 			clientMock: clientServiceMockParams{
// 				ctx: mock.Anything,
// 				client: entity.Client{
// 					ID:       testClientID,
// 					Login:    "sus",
// 					Password: "123",
// 				},
// 				err: nil,
// 			},
// 			sessionMock: sessionRepoMockParams{
// 				ctx:     mock.Anything,
// 				session: mock.Anything,
// 				err:     ErrInternalSessionRepo,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrInternalSessionRepo,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sessionRepo := &mocks.ISessionRepository{}
// 			clientService := &mocks.IClientService{}

// 			sessionRepo.
// 				On("CreateNewSession", tt.sessionMock.ctx, tt.sessionMock.session).
// 				Return(tt.sessionMock.err)
// 			clientService.
// 				On("RegisterNewClient", tt.clientMock.ctx, tt.clientMock.client).
// 				Return(tt.clientMock.err)

// 			a := &AuthorizationService{
// 				sessionRepo:   sessionRepo,
// 				clientService: clientService,
// 			}

// 			got, err := a.Register(tt.args.ctx, tt.args.client)

// 			assert.ErrorIsf(t, err, tt.wantErr, "errors misplace")

// 			assert.Equal(t, tt.want.ClientID, got.ClientID, "wrong client id")
// 		})
// 	}
// }

// func TestAuthorizationService_Logout(t *testing.T) {
// 	type GetSessionBySessionIDParams struct {
// 		ctx       any
// 		sessionID any
// 		session   entity.Session
// 		err       any
// 	}
// 	type DeleteSessionParams struct {
// 		ctx       any
// 		sessionID any
// 		err       any
// 	}
// 	type args struct {
// 		ctx       context.Context
// 		sessionID uuid.UUID
// 	}
// 	tests := []struct {
// 		name              string
// 		args              args
// 		mockGetSession    GetSessionBySessionIDParams
// 		mockDeleteSession DeleteSessionParams
// 		want              entity.Session
// 		wantErr           error
// 	}{
// 		{
// 			name: "logout ok",
// 			args: args{
// 				ctx:       context.Background(),
// 				sessionID: testSessionID,
// 			},
// 			mockDeleteSession: DeleteSessionParams{
// 				ctx:       mock.Anything,
// 				sessionID: testSessionID,
// 				err:       nil,
// 			},
// 			mockGetSession: GetSessionBySessionIDParams{
// 				ctx:       mock.Anything,
// 				sessionID: testSessionID,
// 				session: entity.Session{
// 					ClientID:  testClientID,
// 					SessionID: testSessionID,
// 				},
// 				err: nil,
// 			},
// 			want: entity.Session{
// 				ClientID:  testClientID,
// 				SessionID: testSessionID,
// 			},
// 			wantErr: nil,
// 		},

// 		{
// 			name: "logout fail: error get session",
// 			args: args{
// 				ctx:       context.Background(),
// 				sessionID: testSessionID,
// 			},
// 			mockDeleteSession: DeleteSessionParams{
// 				ctx:       mock.Anything,
// 				sessionID: testSessionID,
// 				err:       nil,
// 			},
// 			mockGetSession: GetSessionBySessionIDParams{
// 				ctx:       mock.Anything,
// 				sessionID: testSessionID,
// 				session:   entity.Session{},
// 				err:       ErrInternalSessionRepo,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrInternalSessionRepo,
// 		},

// 		{
// 			name: "logout fail: error delete session",
// 			args: args{
// 				ctx:       context.Background(),
// 				sessionID: testSessionID,
// 			},
// 			mockDeleteSession: DeleteSessionParams{
// 				ctx:      mock.Anything,
// 				sessionID: testSessionID,
// 				err:      ErrInternalSessionRepo,
// 			},
// 			mockGetSession: GetSessionBySessionIDParams{
// 				ctx:      mock.Anything,
// 				sessionID: testSessionID,
// 				session: entity.Session{
// 					ClientID:  testClientID,
// 					SessionID: testSessionID,
// 				},
// 				err: nil,
// 			},
// 			want:    entity.Session{},
// 			wantErr: ErrInternalSessionRepo,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			sessionRepo := &mocks.ISessionRepository{}
// 			clientService := &mocks.IClientService{}

// 			sessionRepo.
// 				On("GetSessionBySessionID", tt.mockGetSession.ctx, tt.mockGetSession.sessionID).
// 				Return(tt.mockGetSession.session, tt.mockGetSession.err)
// 			sessionRepo.
// 				On("DeleteSession", tt.mockDeleteSession.ctx, tt.mockDeleteSession.sessionID).
// 				Return(tt.mockDeleteSession.err)

// 			a := &AuthorizationService{
// 				sessionRepo:   sessionRepo,
// 				clientService: clientService,
// 			}

// 			got, err := a.Logout(tt.args.ctx, tt.args.sessionID)

// 			assert.ErrorIsf(t, err, tt.wantErr, "errors misplace")

// 			assert.Equal(t, tt.want.ClientID, got.ClientID, "wrong client id")
// 			assert.Equal(t, tt.want.SessionID, got.SessionID, "wrong session id")
// 		})
// 	}
// }
