#define _WINSOCK_DEPRECATED_NO_WARNINGS

#include <winsock2.h>
#include <stdio.h>
#include <windows.h>
#include <ws2tcpip.h>

#pragma comment(lib, "Ws2_32.lib")
//#pragma warning(disable:4996)

void con_msg( USHORT port, PCWSTR ip);
void err(PWSTR code);

DWORD callport;

int CALLBACK ConditionAcceptFunc(
    LPWSABUF lpCallerId,
    LPWSABUF lpCallerData,
    LPQOS pQos,
    LPQOS lpGQOS,
    LPWSABUF lpCalleeId,
    LPWSABUF lpCalleeData,
    GROUP FAR * g,
    DWORD_PTR dwCallbackData
) {
    WSABUF wsab;
    wsab = *(lpCallerId);

    //wprintf_s(L"wsab len: %d\n", wsab.len);
	//wprintf_s(L"wsab buf addr: %#x\n", (DWORD)(WSABUF *)wsab.buf);

	SOCKADDR_STORAGE *piss;
	piss = (SOCKADDR_STORAGE *)wsab.buf;
    if (piss->ss_family == AF_INET) {
        SOCKADDR_IN *meme = (SOCKADDR_IN*)piss;
        IN_ADDR ina = meme->sin_addr;
        wchar_t buf[64];

        wchar_t *pstrbuf = (wchar_t *)SecureZeroMemory( &buf, 64 );
        PWSTR szip = (PWSTR)InetNtopW(AF_INET, &ina, pstrbuf, 64);
        USHORT usport = meme->sin_port;

        //wprintf_s( L"debug ip: %s\n", szip );
        //wprintf_s( L"debug port: %d\n", usport );

        con_msg( callport, szip );
    }

    exit( 0 );

    return CF_DEFER;
}

void rip( ) {
    WSACleanup( );
    exit( 0 );
}

void con_msg( USHORT port, PCWSTR ip) {
    PCWSTR fs = L"{\"returntype\": \"con\", \"port\": %d, \"ip\": \"%s\"}\n";
    wprintf_s(fs, port, ip);
}

void err( PWSTR code ) {
    PCWSTR fs = L"{\"returntype\": \"err\", \"port\": 0, \"ip\": \"%s\"}\n";
    wprintf_s( fs, code );
    exit( 1 );
}

int wmain( int argc, wchar_t *argv[ ], wchar_t *envp[ ] ) {
	setvbuf(stdout, NULL, _IONBF, 0);
	if (argc != 2) {
		err(L"args");
		return 1;
	}

	DWORD port;
	port = _wtoi(argv[1]);

    // validate port no
    //wprintf_s(L"debug port: %d", port);
    if (port < 1 || port > 65535) {
        err(L"portno");
        return 1;
    }
    
    WSADATA wsaData;
    SOCKET ListenSocket, AcceptSocket;
    SOCKADDR_IN saClient;
    int iClientSize = sizeof(saClient);
    char* ip;
    SOCKADDR_IN service;
    int error;

    callport = port;

    error = WSAStartup(MAKEWORD(2,2), &wsaData);
    if (error) {
        err(L"wsastartup");
    }

    ListenSocket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (ListenSocket == INVALID_SOCKET) {
        WSACleanup();
        err(L"socket");
    }

    int enable = 1;

    int rc = setsockopt( 
        ListenSocket,
        SOL_SOCKET,
        SO_CONDITIONAL_ACCEPT,
        &enable,
        sizeof(enable)
    );

    if (rc != 0) {
        err(L"setsockopt");
    }

    service.sin_family = AF_INET;
    service.sin_port = htons(port);
    HOSTENT* thisHost;
    thisHost = gethostbyname("0.0.0.0");
    ip = inet_ntoa (*(struct in_addr *)*thisHost->h_addr_list);
    service.sin_addr.s_addr = inet_addr(ip);

    error = bind(ListenSocket, (SOCKADDR *) &service, sizeof(SOCKADDR));
    if (error == SOCKET_ERROR) {
        closesocket(ListenSocket);
        WSACleanup();
        err(L"bind");
        return 1;
    }

    error = listen(ListenSocket, 1);
    if (error == SOCKET_ERROR) {
        closesocket(ListenSocket);
        WSACleanup();
        err(L"listen");
        return 1;
    }

    //wprintf_s(L"listening");

    //PCWSTR s = L"{\"returntype\": \"start\", \"port\": %d, \"ip\": \"xxx\"}\n";
    //wprintf_s( s, port );


	while (1) {
		AcceptSocket = WSAAccept(ListenSocket, (SOCKADDR*)&saClient, &iClientSize,
			(void *)&ConditionAcceptFunc, NULL);
            rip();
	}

    closesocket(AcceptSocket);
    closesocket(ListenSocket);
    WSACleanup();

    return 0;
}
