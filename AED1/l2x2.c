#include <stdio.h>
#include <stdlib.h>

struct Nota{
    int Num;
    int Data;
    int Resp;
    float Valor;
}nota;

typedef no no;

struct no{
    struct Nota Dado;
    struct no* Pai;
    struct no* Dir;
    struct no* Esq;
};

no *raiz;

struct{
    int Dias;
    int Qntfunc;
}reg;

int contdias(int mes);
void iniciar();

int temp;
float ftemp;
void main(){

    {   /*Leitura do reg*/
    do{
        printf("Entre com o mes em que as vendas foram realizadas: ");
        scanf("%i",&temp);printf("\n\n");

        if(temp>12 || temp<1){
            printf("Mes deve ser maior ou igual a 1 e menor ou igual a 12\n\n");
        }
    }while(temp>12 || temp<1);
    reg.Dias = contdias(temp); /*guarda os dias*/

    do{
        printf("Entre com o numero de vendedores que trabalharam no mes: ");
        scanf("%i",&temp);printf("\n\n");

        if(temp<=0){
            printf("Numero de vendedores deve ser maior do que zero\n\n");
        }
    }while(temp<=0);
    reg.Qntfunc = temp; /*guarda os funcionarios*/

    }

    {/*leitura da venda*/
    int flag,cont;
    printf("Entre com as vendas. Para cada venda realizada no mes informar o numero da nota fiscal que a identifica, o codigo do vendedor responsavel, o dia do mes em que foi feita e o valor da venda em reais. Quando desejar finalizar a entrada de dados digite o numero zero para a nota fiscal.\n\n");
    no* raiz;
    do{
        flag = 0;
        printf("Venda: ");
        scanf("%i",&temp);
        if (temp==0){
            break;
        }
        nota.Num = temp;

        scanf("%i",&temp);
        if (temp>reg.Qntfunc){
            printf("Codigo do vendedor invalido\n\n");
            flag = 1; 
        }
        nota.Resp = temp;

        scanf("%i",&temp);
        if (temp>reg.Dias || temp<1){
            printf("Dia da venda invalido\n\n");
            flag = 1;
        }
        nota.Data = temp;

        scanf("%f",&ftemp);
        if (ftemp<=0){
            printf("Valor da venda invalida\n\n");
            flag = 1;
        }
        nota.Valor = ftemp;
        if (flag==0){
            inserir();
        }
    }while(1);

    }
}


int contdias(int mes){
    if(mes == 2){
        return 28;
    }else if (mes == 1 || mes == 3 || mes == 5 || mes == 7 || mes == 8 || mes == 10 || mes == 12){
        return 31;
    }else{
        return 30;
    }
}


void iniciar(){
    raiz = NULL;
}

void inserir(no *traiz,struct Nota dado){
	
    if(traiz==NULL){
        traiz = malloc(sizeof(no));
        traiz->Esq = NULL;
        traiz->Dir = NULL;
        traiz->Pai = NULL;
    }
}


