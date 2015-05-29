import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.io.*;
import java.util.Date;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Calendar;

import com.rti.dds.domain.*;
import com.rti.dds.infrastructure.*;
import com.rti.dds.publication.*;
import com.rti.dds.topic.*;
import com.rti.dds.typecode.*;

public class TopFilter implements ContentFilter
{
	private interface EvaluateFunction {
        public boolean eval(float x, float y);
    };
	

    private class ValidityTest implements EvaluateFunction
	{
		float _cpuuse = 0.02f;
		float _memuse = 0.50f;

        public ValidityTest(float cpu, float mem) {
            _cpuuse = cpu;
			_memuse = mem;
        }

        public boolean eval(float samplecpu, float samplemem) {
			if (samplecpu >= _cpuuse && samplemem >= _memuse)
				return true;
			return false;
        }
    };
	
    public void compile(ObjectHolder new_compile_data, String expression,
            StringSeq parameters, TypeCode type_code, String type_class_name,
            Object old_compile_data) {
        /*
         * We expect an expression of the form "%0 %1 <var>" where %1 =
         * "divides" or "greater-than" and <var> is an integral member of the
         * msg structure.
         * 
         * We don't actually check that <var> has the correct typecode, (or even
         * that it exists!). See example Typecodes to see how to work with
         * typecodes.
         * 
         * The compile information is a structure including the first filter
         * parameter (%0) and a function pointer to evaluate the sample
         */

        // Check form:
		System.out.println("Expression: " + expression);
		System.out.print("Parameters: ");
		for (int i = 0; i < parameters.size(); i++)
		{
			System.out.print((String) parameters.get(i) + ", ");
		}
		System.out.println("~END~");
		
        if (expression.startsWith("%0 %1 %2 ") && expression.length() > 4
                && parameters.size() > 2) { // Enough parameters?

            float cpu = Float.valueOf((String) parameters.get(0)).floatValue();
            float mem = Float.valueOf((String) parameters.get(1)).floatValue();

            if (parameters.get(3).equals("valid")) {
                new_compile_data.value = new ValidityTest(cpu, mem);
                return;
            }
        }

        System.out.print("CustomFilter: Unable to compile expression '"
                + expression + "'\n");
        System.out.print("              with parameters '" + parameters.get(0)
                + "' '" + parameters.get(1)+ "' '" + parameters.get(2) + "'\n");
        throw (new RETCODE_BAD_PARAMETER());
    }

    public boolean evaluate(Object compile_data, Object sample,
            FilterSampleInfo meta_data) {
        float x = ((TopFunction) sample).cpuUsage;
        float y = ((TopFunction) sample).memUsage;
		
        return ((EvaluateFunction) compile_data).eval(x, y);
    }

    public void finalize(Object compile_data) {
    }
}

