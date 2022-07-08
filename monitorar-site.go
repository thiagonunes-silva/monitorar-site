package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const STATUS_CODE_SUCESSO = 200
const NOME_ARQUIVO_SITES = "sites.txt"
const ARQUIVO_LOG_MONITORACAO = "log.txt"

func main() {
	opcaoSelecionada := 0
	exibirSaudacao()
	for opcaoSelecionada != 4 {
		exibirMenu()
		opcaoSelecionada = lerOpcaoUsuario()
		tratarOpcaoSelecionada(opcaoSelecionada)
	}
}

func tratarOpcaoSelecionada(opcaoSelecionada int) {
	switch opcaoSelecionada {
	case 1:
		iniciarMonitoramento()
	case 2:
		exibirLogs()
	case 3:
		registrarNovoSiteParaMonitoracao()
	case 4:
		exibirMensagemDeSaida()
		os.Exit(0)
	default:
		fmt.Println("\nOpção INVÁLIDA!!!!!!!!!!!!")
		//os.Exit(0)
	}

}

func lerOpcaoUsuario() int {
	opcaoSelecionada := 0
	fmt.Scan(&opcaoSelecionada)
	return opcaoSelecionada
}

func iniciarMonitoramento() {
	fmt.Println("Iniciando Monitoramento!")
	listaSites := obterListaSitesParaMonitoracao()
	for _, site := range listaSites {
		monitorarSite(site)
	}
}

func monitorarSite(site string) {
	if site != "" {
		resp, err := http.Get(site)
		if err != nil {
			fmt.Println("Erro ao tentar acesso o site", site, ". Erro:", err)
		} else {
			if resp.StatusCode == STATUS_CODE_SUCESSO {
				fmt.Printf("Site %s foi carregado com sucesso! StatusCode %d \n", site, resp.StatusCode)
			} else {
				fmt.Printf("Site %s não foi carregado com sucesso! StatusCode %d \n", site, resp.StatusCode)
			}
			registrarLog(site, resp.StatusCode)
		}
	}
}

func obterListaSitesParaMonitoracao() []string {
	var listaSites []string
	arquivo, err := os.Open(NOME_ARQUIVO_SITES)
	defer arquivo.Close()
	if err != nil {
		fmt.Println("Erro ao tentar abrir o arquivo", NOME_ARQUIVO_SITES, "Erro: ", err)
	} else {
		leitor := bufio.NewReader(arquivo)
		for {
			linhaArquivo, err := leitor.ReadString('\n')
			listaSites = append(listaSites, strings.TrimSpace(linhaArquivo))
			if err == io.EOF {
				break
			}
		}
	}
	return listaSites
}

func exibirLogs() {
	arquivo, err := ioutil.ReadFile(ARQUIVO_LOG_MONITORACAO)
	if err != nil {
		fmt.Println("Erro ao tentar abrir o arquivo", ARQUIVO_LOG_MONITORACAO, "Erro:", err)
	} else {
		fmt.Println(string(arquivo))
	}
}

func exibirSaudacao() {
	fmt.Println("\nOlá, seja bem vindo!")
}

func exibirMenu() {
	fmt.Println("\nSelecione uma das opções abaixo:")
	fmt.Println("\n\t1 - Executar monitoramento")
	fmt.Println("\t2 - Exibir Logs")
	fmt.Println("\t3 - Registrar Site para Monitoração")
	fmt.Println("\t4 - Sair do programa")
}

func exibirMensagemDeSaida() {
	fmt.Println("\nObrigado, até a próxima!")
}

func registrarLog(site string, statusCode int) {
	arquivoLog, err := os.OpenFile(ARQUIVO_LOG_MONITORACAO, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer arquivoLog.Close()
	if err != nil {
		fmt.Println("Erro ao tentar abrir o arquivo", ARQUIVO_LOG_MONITORACAO, "Erro:", err)
	} else {
		arquivoLog.WriteString("Monitoracao[" + time.Now().Format("02/01/2006 15:04:05") + "] - Site[" + site + "] - StatusCode=[" + strconv.FormatInt(int64(statusCode), 10) + "]\n")
	}
}

func registrarNovoSiteParaMonitoracao() {
	fmt.Print("Informe o novo site para monitoração: ")
	var novoSite string
	fmt.Scan(&novoSite)
	listaSitesCadastrados := obterListaSitesParaMonitoracao()
	for _, siteCadastrado := range listaSitesCadastrados {
		if siteCadastrado == novoSite {
			fmt.Println("O Site informado já está cadastrado")
			return
		}
	}
	registrarSiteParaMonitoracao(novoSite)
	fmt.Println("Site registrado com sucesso")
}

func registrarSiteParaMonitoracao(site string) {
	arquivoSites, err := os.OpenFile(NOME_ARQUIVO_SITES, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer arquivoSites.Close()
	if err != nil {
		fmt.Println("Erro ao tentar abrir o arquivo", NOME_ARQUIVO_SITES, "Erro:", err)
	} else {
		arquivoSites.WriteString(site + "\n")
	}
}
