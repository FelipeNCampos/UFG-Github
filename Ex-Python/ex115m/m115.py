# opc1 opc2 opc3
# escolha
def opc1(msg1):
    print('---'*15)
    print('Opção 1')
    print('---'*15)



def opc2(msg2):
    print('---'*15)
    print('Opção 2')
    print('---'*15)




def opc3(msg3):
    print('---'*15)
    print('Opção 3')
    print('---'*15)


def choisse(msg):
    while True:
        try :
            op = int(input(msg))
        except TypeError:
            print('\033[031mERRO\033[m,por favor insira um valor inteiro válido.')
        except:
            print('\033[031mERRO\033[m, tente novamente.')
        else :
            if op==1:
                try:
                    print('---'*15)
                    arq = open('Registro.txt', 'r')
                    print(arq.read())
                    print('---'*15)
                except:
                    print('\033[031mERRO\033[m ao ler o registro, tente novamente mais tarde.')

            elif op==2:
                arq = open('Registro.txt', 'a')
                while True:
                    dados = {}
                    dados['Nome'] = input('Digite o nome a ser registrado :')
                    dados['Idade'] = input(f'Digite a idade do indivíduo {dados['Nome']} : ')
                    imp = ''
                    for x,c in enumerate(dados):
                        if x>=1:
                            imp += '|'
                        imp += c+ '='+ dados[c]
                    print(imp)
                    arq.write(imp)
                    dados.clear()
                    try :
                        rrr = str(input('Deseja continuar (S/N) ?')).strip().upper()
                        if rrr=='N':
                            break
                        if rrr not in 'NS':
                            print('\033[031mResposta inválida\033[m, vou tomar como um não.')
                        if rrr=='S':
                            continue
                    except:
                        print('\033[31mResposta inválida, vou tomar como um nâo.')
                        break
            elif op==3:
                print('-=-'*15)
                print('\033[031mENCERRANDO PROGRAMA\033[m')
                return
            elif op>3 or op<=0:
                print('\033[031mERRO\033[m,Opção fora do limite.')

